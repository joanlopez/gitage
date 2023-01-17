//go:build !windows

package fstest

func rootDir() string {
	return pathSepStr
}
