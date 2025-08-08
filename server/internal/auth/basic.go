package auth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"asciinema-server-x/server/internal/util"
)

// AdminBasic requires admin Basic Auth for protected endpoints.
func AdminBasic(cfg util.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, p, ok := parseBasic(c.Request.Header.Get("Authorization"))
		if !ok || u != cfg.AdminUser || p != cfg.AdminPass {
			util.Error(c, http.StatusUnauthorized, "INVALID_AUTH", "unauthorized")
			return
		}
		c.Next()
	}
}

// ParseBasic returns username & password from Authorization header.
func ParseBasic(h string) (string, string, bool) { return parseBasic(h) }

func parseBasic(h string) (string, string, bool) {
	if h == "" {
		return "", "", false
	}
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Basic") {
		return "", "", false
	}
	b, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", false
	}
	up := string(b)
	i := strings.IndexByte(up, ':')
	if i < 0 {
		return "", "", false
	}
	return up[:i], up[i+1:], true
}
