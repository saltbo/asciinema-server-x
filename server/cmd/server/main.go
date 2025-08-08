package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"asciinema-server-x/server/internal/auth"
	"asciinema-server-x/server/internal/handler"
	"asciinema-server-x/server/internal/storage"
	"asciinema-server-x/server/internal/util"
)

func main() {
	cfg := util.LoadConfig()

	// Ensure storage root exists
	if err := storage.EnsureDir(cfg.StorageRoot); err != nil {
		log.Fatalf("failed to ensure storage root: %v", err)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Healthz
	r.GET("/healthz", handler.Healthz())

	// Public upload (uses username/machine-id Basic per asciinema CLI)
	r.POST("/api/asciicasts", handler.UploadCast(cfg))

	// Admin-protected routes
	api := r.Group("/api", auth.AdminBasic(cfg))
	{
		api.POST("/users", handler.CreateOrGetUser(cfg))
		api.GET("/users", handler.ListUsers(cfg))
		api.GET("/users/:username/casts", handler.ListUserCasts(cfg))
		api.GET("/casts/file", handler.GetCastFile(cfg))
	}

	// Static files for SPA if present
	if st, err := os.Stat(cfg.WebDistRoot); err == nil && st.IsDir() {
		r.Static("/assets", filepath.Join(cfg.WebDistRoot, "assets"))
		r.StaticFile("/favicon.ico", filepath.Join(cfg.WebDistRoot, "favicon.ico"))
		r.NoRoute(func(c *gin.Context) {
			// For API 404, keep JSON
			if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
				c.Status(http.StatusNotFound)
				return
			}
			c.File(filepath.Join(cfg.WebDistRoot, "index.html"))
		})
	}

	addr := ":" + cfg.HTTPPort
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
