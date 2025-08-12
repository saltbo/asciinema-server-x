package handler

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"asciinema-server-x/server/internal/auth"
	"asciinema-server-x/server/internal/storage"
	"asciinema-server-x/server/internal/util"
)

type uploadResp struct {
	Username string `json:"username"`
	RelPath  string `json:"relPath"`
	Size     int64  `json:"sizeBytes"`
	URL      string `json:"url"`
}

// CastMetadata represents the header of an asciinema cast file
type CastMetadata struct {
	Version   int     `json:"version"`
	Width     int     `json:"width"`
	Height    int     `json:"height"`
	Timestamp int64   `json:"timestamp"`
	Title     string  `json:"title,omitempty"`
	Duration  float64 `json:"duration,omitempty"`
}

type castItem struct {
	RelPath  string        `json:"relPath"`
	Size     int64         `json:"sizeBytes"`
	MTime    time.Time     `json:"mtime"`
	Metadata *CastMetadata `json:"metadata,omitempty"`
}

type castListResp struct {
	Items []castItem `json:"items"`
	Total int        `json:"total"`
}

// parseCastMetadata reads the first line of a cast file and parses the metadata
func parseCastMetadata(filePath string) (*CastMetadata, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}

	var metadata CastMetadata
	if err := json.Unmarshal(scanner.Bytes(), &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// encodeMetadataToFilename creates a compact filename with encoded metadata
func encodeMetadataToFilename(metadata *CastMetadata) string {
	// Use metadata timestamp if available, otherwise generate current timestamp
	timestamp := metadata.Timestamp
	if timestamp == 0 {
		timestamp = time.Now().Unix()
	}

	// Create a compact JSON representation of the metadata including timestamp
	metaData := map[string]interface{}{
		"w":  metadata.Width,
		"h":  metadata.Height,
		"t":  metadata.Title,
		"ts": timestamp,
	}

	// Add duration if available
	if metadata.Duration > 0 {
		metaData["d"] = metadata.Duration
	}

	metaJSON, _ := json.Marshal(metaData)
	// Use base64 URL encoding (safe for filenames) - add .cast extension for filesystem
	metaEncoded := base64.RawURLEncoding.EncodeToString(metaJSON)

	return metaEncoded + ".cast"
}

// decodeMetadataFromFilename extracts metadata from encoded filename
func decodeMetadataFromFilename(filename string) *CastMetadata {
	// Remove .cast extension
	encoded := strings.TrimSuffix(filename, ".cast")

	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil
	}

	var metaData map[string]interface{}
	if err := json.Unmarshal(decoded, &metaData); err != nil {
		return nil
	}

	metadata := &CastMetadata{}
	if w, ok := metaData["w"].(float64); ok {
		metadata.Width = int(w)
	}
	if h, ok := metaData["h"].(float64); ok {
		metadata.Height = int(h)
	}
	if t, ok := metaData["t"].(string); ok {
		metadata.Title = t
	}
	if ts, ok := metaData["ts"].(float64); ok {
		metadata.Timestamp = int64(ts)
	}
	if d, ok := metaData["d"].(float64); ok {
		metadata.Duration = d
	}

	return metadata
}

// calculateDuration calculates duration by parsing the cast file to find the last timestamp
func calculateDuration(filePath string, metadata *CastMetadata) {
	// If duration is already present, don't calculate
	if metadata.Duration > 0 {
		return
	}

	// Parse the cast file to find the last timestamp
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lastTimestamp float64 = 0

	// Skip the first line (header)
	if !scanner.Scan() {
		return
	}

	// Read through all lines to find the last timestamp
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		// Parse the line as [timestamp, event_type, data]
		var entry []interface{}
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		if len(entry) >= 1 {
			if timestamp, ok := entry[0].(float64); ok {
				if timestamp > lastTimestamp {
					lastTimestamp = timestamp
				}
			}
		}
	}

	// Set the duration to the last timestamp (which represents total recording time)
	if lastTimestamp > 0 {
		metadata.Duration = lastTimestamp
	}
}

func UploadCast(cfg util.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse Basic: username:machine-id
		u, p, ok := auth.ParseBasic(c.Request.Header.Get("Authorization"))
		if !ok || !validUsername(u) {
			util.Error(c, http.StatusUnauthorized, "INVALID_AUTH", "invalid basic auth")
			return
		}
		mid, err := storage.ReadUserMachineID(cfg.StorageRoot, u)
		if err != nil || p != strings.TrimSpace(mid) {
			util.Error(c, http.StatusUnauthorized, "INVALID_AUTH", "machine-id mismatch")
			return
		}
		file, err := c.FormFile("asciicast")
		if err != nil {
			util.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "missing asciicast field")
			return
		}
		if cfg.MaxUploadMB > 0 && file.Size > cfg.MaxUploadMB*1024*1024 {
			util.Error(c, http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE", "file too large")
			return
		}
		f, err := file.Open()
		if err != nil {
			util.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}
		defer f.Close()

		// Read the first line to extract metadata
		scanner := bufio.NewScanner(f)
		if !scanner.Scan() {
			util.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "invalid cast file format")
			return
		}

		var metadata CastMetadata
		if err := json.Unmarshal(scanner.Bytes(), &metadata); err != nil {
			util.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "invalid cast file header")
			return
		}

		// Reset file position
		f.Close()
		f, err = file.Open()
		if err != nil {
			util.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}
		defer f.Close()

		date := util.TodayYYYYMMDD()
		userDir, err := storage.EnsureUser(cfg.StorageRoot, u)
		if err != nil {
			util.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}
		dateDir := filepath.Join(userDir, date)
		if err := os.MkdirAll(dateDir, 0o755); err != nil {
			util.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}

		// Create filename with encoded metadata (no UUID needed)
		filename := encodeMetadataToFilename(&metadata)
		castPath := filepath.Join(dateDir, filename)
		n, err := storage.WriteFileAtomic(castPath, f)
		if err != nil {
			util.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}
		rel := filepath.Join(u, date, filename)
		// Build absolute URL for the frontend Player page
		var base string
		if cfg.PublicBaseURL != "" {
			base = strings.TrimRight(cfg.PublicBaseURL, "/")
		} else {
			// Derive from request (honor common proxy headers)
			scheme := c.Request.Header.Get("X-Forwarded-Proto")
			if scheme == "" {
				if c.Request.TLS != nil {
					scheme = "https"
				} else {
					scheme = "http"
				}
			}
			host := c.Request.Header.Get("X-Forwarded-Host")
			if host == "" {
				host = c.Request.Host
			}
			base = scheme + "://" + host
			// Optional path prefix if behind a reverse proxy: X-Forwarded-Prefix
			if p := strings.TrimSuffix(c.Request.Header.Get("X-Forwarded-Prefix"), "/"); p != "" {
				base += p
			}
		}
		playURL := base + "/play/" + url.PathEscape(rel)
		c.JSON(http.StatusCreated, uploadResp{Username: u, RelPath: rel, Size: n, URL: playURL})
	}
}

func ListUserCasts(cfg util.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := c.Param("username")
		if !validUsername(u) {
			util.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "invalid username")
			return
		}
		root := filepath.Join(cfg.StorageRoot, u)
		if st, err := os.Stat(root); err != nil || !st.IsDir() {
			util.Error(c, http.StatusNotFound, "USER_NOT_FOUND", "user not found")
			return
		}
		entries, _ := os.ReadDir(root)
		var items []castItem
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			dd := filepath.Join(root, e.Name())
			files, _ := os.ReadDir(dd)
			for _, f := range files {
				name := f.Name()
				// Check for .cast extension
				if !strings.HasSuffix(name, ".cast") {
					continue
				}

				full := filepath.Join(dd, name)
				st, err := os.Stat(full)
				if err != nil {
					continue
				}

				// Remove .cast extension from relPath for API response
				nameWithoutExt := strings.TrimSuffix(name, ".cast")
				rel := filepath.Join(u, e.Name(), nameWithoutExt)

				// Try to decode metadata from filename first
				var metadata *CastMetadata
				if decodedMetadata := decodeMetadataFromFilename(name); decodedMetadata != nil {
					metadata = decodedMetadata
				} else {
					// Fallback: parse metadata from file content
					if castMetadata, err := parseCastMetadata(full); err == nil {
						metadata = castMetadata
					}
				}

				// Calculate duration if not present
				if metadata != nil {
					calculateDuration(full, metadata)
				}

				items = append(items, castItem{
					RelPath:  rel,
					Size:     st.Size(),
					MTime:    st.ModTime(),
					Metadata: metadata,
				})
			}
		}
		c.JSON(http.StatusOK, castListResp{Items: items, Total: len(items)})
	}
}

func GetCastFile(cfg util.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		rel := c.Query("path")
		// Add .cast extension for filesystem lookup
		relWithExt := rel + ".cast"
		full, err := storage.SafeJoin(cfg.StorageRoot, relWithExt)
		if err != nil {
			util.Error(c, http.StatusBadRequest, "PATH_TRAVERSAL_BLOCKED", "invalid path")
			return
		}
		if fi, err := os.Stat(full); err != nil || fi.IsDir() {
			util.Error(c, http.StatusNotFound, "CAST_NOT_FOUND", "not found")
			return
		}
		// Stream file
		f, err := os.Open(full)
		if err != nil {
			util.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}
		defer f.Close()
		c.Header("Content-Type", "application/octet-stream")
		c.Status(http.StatusOK)
		_, _ = io.Copy(c.Writer, f)
	}
}
