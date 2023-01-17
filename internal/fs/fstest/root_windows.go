//go:build windows

package fstest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Rootify(path string) string {
	path = strings.ReplaceAll(path[1:], "/", pathSepStr)
	return filepath.Clean(fmt.Sprintf("%s%s%s", rootDir(), pathSepStr, path))
}

func rootDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s%s", strings.Split(cwd, pathSepStr)[0], pathSepStr)
}
