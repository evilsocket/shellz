package core

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	defaultTrimSet = "\r\n\t "
)

func Trim(s string) string {
	return strings.Trim(s, defaultTrimSet)
}

func TrimRight(s string) string {
	return strings.TrimRight(s, defaultTrimSet)
}

func TrimLeft(s string) string {
	return strings.TrimLeft(s, defaultTrimSet)
}

func SepSplit(sv string, sep string) []string {
	filtered := make([]string, 0)
	for _, part := range strings.Split(sv, sep) {
		part = Trim(part)
		if part != "" {
			filtered = append(filtered, part)
		}
	}
	return filtered
}

func CommaSplit(csv string) []string {
	return SepSplit(csv, ",")
}

func ExpandPath(path string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Abs(strings.Replace(path, "~", usr.HomeDir, 1))
}

func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func Glob(path string, expr string, cb func(fileName string) error) error {
	if files, err := filepath.Glob(filepath.Join(path, expr)); err != nil {
		return err
	} else {
		for _, fileName := range files {
			if err := cb(fileName); err != nil {
				return err
			}
		}
	}
	return nil
}
