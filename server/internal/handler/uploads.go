package handler

import (
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

type castItem struct {
	RelPath string    `json:"relPath"`
	Size    int64     `json:"sizeBytes"`
	MTime   time.Time `json:"mtime"`
}

type castListResp struct {
	Items []castItem `json:"items"`
	Total int        `json:"total"`
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
		id, err := util.NewUUID()
		if err != nil {
			util.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}
		castPath := filepath.Join(dateDir, id+".cast")
		n, err := storage.WriteFileAtomic(castPath, f)
		if err != nil {
			util.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}
		rel := filepath.Join(u, date, id+".cast")
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
				if !strings.HasSuffix(name, ".cast") {
					continue
				}
				full := filepath.Join(dd, name)
				st, err := os.Stat(full)
				if err != nil {
					continue
				}
				rel := filepath.Join(u, e.Name(), name)
				items = append(items, castItem{RelPath: rel, Size: st.Size(), MTime: st.ModTime()})
			}
		}
		c.JSON(http.StatusOK, castListResp{Items: items, Total: len(items)})
	}
}

func GetCastFile(cfg util.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		rel := c.Query("path")
		full, err := storage.SafeJoin(cfg.StorageRoot, rel)
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
