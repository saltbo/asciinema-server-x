package storage

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func EnsureDir(p string) error { return os.MkdirAll(p, 0o755) }

func EnsureUser(root, username string) (string, error) {
	userDir := filepath.Join(root, username)
	if err := os.MkdirAll(userDir, 0o755); err != nil {
		return "", err
	}
	return userDir, nil
}

func WriteFileAtomic(dstPath string, r io.Reader) (int64, error) {
	tmp := dstPath + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	n, err := io.Copy(f, r)
	if err != nil {
		_ = os.Remove(tmp)
		return 0, err
	}
	if err := f.Sync(); err != nil {
		_ = os.Remove(tmp)
		return 0, err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return 0, err
	}
	if err := os.Rename(tmp, dstPath); err != nil {
		_ = os.Remove(tmp)
		return 0, err
	}
	return n, nil
}

func ReadUserMachineID(root, username string) (string, error) {
	b, err := os.ReadFile(filepath.Join(root, username, "machine-id"))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

func EnsureUserMachineID(root, username, machineID string) error {
	return os.WriteFile(filepath.Join(root, username, "machine-id"), []byte(machineID+"\n"), 0o644)
}

func SafeJoin(root, rel string) (string, error) {
	clean := filepath.Clean(rel)
	if strings.HasPrefix(clean, "..") || filepath.IsAbs(clean) {
		return "", errors.New("invalid path")
	}
	full := filepath.Join(root, clean)
	relchk, err := filepath.Rel(root, full)
	if err != nil || strings.HasPrefix(relchk, "..") {
		return "", errors.New("invalid path")
	}
	return full, nil
}
