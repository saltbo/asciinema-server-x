package util

import (
	"os"
)

type Config struct {
	StorageRoot string
	AdminUser   string
	AdminPass   string
	HTTPPort    string
	MaxUploadMB int64
	WebDistRoot string
	// PublicBaseURL is optional; if set, used to build absolute URLs in API responses.
	PublicBaseURL string
}

func LoadConfig() Config {
	cfg := Config{
		StorageRoot:   getenv("STORAGE_ROOT", "./data"),
		AdminUser:     getenv("ADMIN_BASIC_USER", "admin"),
		AdminPass:     getenv("ADMIN_BASIC_PASS", "admin"),
		HTTPPort:      getenv("HTTP_PORT", "8080"),
		MaxUploadMB:   50,
		WebDistRoot:   getenv("WEB_DIST_ROOT", "../web/dist"),
		PublicBaseURL: getenv("PUBLIC_BASE_URL", ""),
	}
	return cfg
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
