package handler

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"asciinema-server-x/server/internal/auth"
	"asciinema-server-x/server/internal/storage"
	"asciinema-server-x/server/internal/util"
)

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
	ShortID  string        `json:"shortId"`
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

		// Generate a cast ID
		castID, err := util.GenerateCastID(u, date)
		if err != nil {
			util.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}
		filename := castID.UniqueID + ".cast"
		castPath := filepath.Join(dateDir, filename)
		if _, err := storage.WriteFileAtomic(castPath, f); err != nil {
			util.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}

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

		c.JSON(http.StatusCreated, gin.H{"url": base + "/a/" + castID.Encode()})
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

				// Remove .cast extension from filename to get the uniqueID
				uniqueID := strings.TrimSuffix(name, ".cast")

				// Create cast ID
				castID := &util.CastID{
					Username: u,
					Date:     e.Name(),
					UniqueID: uniqueID,
				}

				// Parse metadata from file content
				var metadata *CastMetadata
				if castMetadata, err := parseCastMetadata(full); err == nil {
					metadata = castMetadata
					// Calculate duration if not present
					calculateDuration(full, metadata)
				}

				items = append(items, castItem{
					ShortID:  castID.Encode(),
					Size:     st.Size(),
					MTime:    st.ModTime(),
					Metadata: metadata,
				})
			}
		}
		c.JSON(http.StatusOK, castListResp{Items: items, Total: len(items)})
	}
}

// GetCastMetaByID handles requests to get cast metadata by encoded ID
func GetCastMetaByID(cfg util.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			util.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "missing id parameter")
			return
		}

		// Decode ID to get cast information
		castID, err := util.DecodeCastID(id)
		if err != nil {
			util.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid id format")
			return
		}

		// Build file path
		filename := castID.UniqueID + ".cast"
		full := filepath.Join(cfg.StorageRoot, castID.Username, castID.Date, filename)

		fi, err := os.Stat(full)
		if err != nil || fi.IsDir() {
			util.Error(c, http.StatusNotFound, "CAST_NOT_FOUND", "not found")
			return
		}

		// Parse metadata from file content
		var metadata *CastMetadata
		if castMetadata, err := parseCastMetadata(full); err == nil {
			metadata = castMetadata
			// Calculate duration if not present
			calculateDuration(full, metadata)
		}

		// Return the same shortID that was passed in
		item := castItem{
			ShortID:  id,
			Size:     fi.Size(),
			MTime:    fi.ModTime(),
			Metadata: metadata,
		}

		c.JSON(http.StatusOK, item)
	}
}

// GetCastFileByID handles requests to get cast file content by encoded ID
func GetCastFileByID(cfg util.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			util.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "missing id parameter")
			return
		}

		// Decode ID to get cast information
		castID, err := util.DecodeCastID(id)
		if err != nil {
			util.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid id format")
			return
		}

		// Build file path
		filename := castID.UniqueID + ".cast"
		full := filepath.Join(cfg.StorageRoot, castID.Username, castID.Date, filename)

		if fi, err := os.Stat(full); err != nil || fi.IsDir() {
			util.Error(c, http.StatusNotFound, "CAST_NOT_FOUND", "not found")
			return
		}

		// Stream file content
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
