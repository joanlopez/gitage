//go:build !windows

package fs

// implementing normalizeLongPath as stub prevents too much duplicated code in os_windows.go
func normalizeLongPath(path string) string {
	return path
}
