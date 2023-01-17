//go:build !windows

package fstest

func Rootify(path string) string {
	return path
}

func rootDir() string {
	return pathSepStr
}
