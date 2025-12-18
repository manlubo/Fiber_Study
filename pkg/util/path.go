package util

import (
	"os"
	"path/filepath"
)

// 프로젝트 루트 경로
func ProjectRoot() string {
	cwd, _ := os.Getwd()
	return filepath.Dir(filepath.Dir(cwd))
}

// 패스를 루트에서 부터 조립
func GetPath(path string) string {
	return filepath.Join(ProjectRoot(), path)
}
