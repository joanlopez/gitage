package gitage

import (
	"path/filepath"
)

func dir(path string) string {
	return filepath.Join(path, ".gitage")
}
