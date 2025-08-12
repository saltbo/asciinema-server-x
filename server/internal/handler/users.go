package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"asciinema-server-x/server/internal/storage"
	"asciinema-server-x/server/internal/util"
)

type userReq struct {
	Username string `json:"username"`
}

type createResp struct {
	Username  string `json:"username"`
	MachineID string `json:"machineId"`
}

type userListResp struct {
	Items []string `json:"items"`
	Total int      `json:"total"`
}

func CreateOrGetUser(cfg util.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var r userReq
		if err := c.BindJSON(&r); err != nil || !validUsername(r.Username) {
			util.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "invalid username")
			return
		}
		userDir, err := storage.EnsureUser(cfg.StorageRoot, r.Username)
		if err != nil {
			util.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}
		midPath := filepath.Join(userDir, "machine-id")
		if b, err := os.ReadFile(midPath); err == nil {
			c.JSON(http.StatusOK, createResp{Username: r.Username, MachineID: strings.TrimSpace(string(b))})
			return
		}
		mid := uuid.New().String()
		if err := os.WriteFile(midPath, []byte(mid+"\n"), 0o644); err != nil {
			util.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}
		c.JSON(http.StatusCreated, createResp{Username: r.Username, MachineID: mid})
	}
}

func ListUsers(cfg util.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		dirs, err := os.ReadDir(cfg.StorageRoot)
		if err != nil {
			util.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}
		var users []string
		for _, e := range dirs {
			if !e.IsDir() {
				continue
			}
			name := e.Name()
			if !validUsername(name) {
				continue
			}
			if _, err := os.Stat(filepath.Join(cfg.StorageRoot, name, "machine-id")); err == nil {
				users = append(users, name)
			}
		}
		sort.Strings(users)
		c.JSON(http.StatusOK, userListResp{Items: users, Total: len(users)})
	}
}

func validUsername(u string) bool {
	if u == "" || strings.HasPrefix(u, ".") {
		return false
	}
	for _, r := range u {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			continue
		}
		return false
	}
	return true
}
