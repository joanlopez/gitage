//go:build windows

package fs

import (
	"strings"
)

// windows cannot handle long relative paths, so convert to absolute paths by default
func normalizeLongPath(path string) string {
	// path is already normalized
	if strings.HasPrefix(absolutePath, `\\?\`) {
		return absolutePath
	}

	// normalize "network path"
	if strings.HasPrefix(absolutePath, `\\`) {
		return `\\?\UNC\` + strings.TrimPrefix(absolutePath, `\\`)
	}

	// normalize "local path"
	return `\\?\` + absolutePath
}
