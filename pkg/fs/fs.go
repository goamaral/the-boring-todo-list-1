package fs

import (
	"path"
	"path/filepath"
	"runtime"
)

func RelativePath(relativePath string) string {
	_, file, _, _ := runtime.Caller(1)
	return path.Join(filepath.Dir(file), relativePath)
}
