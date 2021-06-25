package xutils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func init() {
	initRelative()
}

var prefixPath string

func initRelative() {
	_, fileName, _, _ := runtime.Caller(0)
	prefixPath = filepath.ToSlash(filepath.Dir(filepath.Dir(filepath.Dir(fileName)))) + "/"
}

func GetPwd() string {
	_, fileName, _, _ := runtime.Caller(1)
	return filepath.ToSlash(filepath.Dir(fileName))
}

func AbsPath(path string) string {
	return prefixPath + filepath.ToSlash(path)
}

func RelativePath(path string) string {
	return strings.TrimPrefix(filepath.ToSlash(path), prefixPath)
}

func FileExists(filename string) bool {
	if filename == "" {
		return false
	}
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
