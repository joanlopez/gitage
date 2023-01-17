//go:build windows

package fstest

import (
	"os"
)

func rootDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return strings.Split(cwd, pathSepStr)[0]
}
