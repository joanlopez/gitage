//go:build windows

package fs

import (
	"strings"
)

// windows cannot handle long relative paths, so convert to absolute paths by default
func normalizeLongPath(path string) string {
	// path is already normalized
	if strings.HasPrefix(path, `\\?\`) {
		return path
	}

	// normalize "network path"
	if strings.HasPrefix(path, `\\`) {
		return `\\?\UNC\` + strings.TrimPrefix(path, `\\`)
	}

	// normalize "local path"
	return `\\?\` + path
}
