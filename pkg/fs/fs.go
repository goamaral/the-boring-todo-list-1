package fs

import (
	"path/filepath"
	"runtime"
)

func RelativePath(relativePath string) string {
	_, file, _, _ := runtime.Caller(1)
	folderPath := filepath.Dir(file)
	return folderPath + "/" + relativePath
}
