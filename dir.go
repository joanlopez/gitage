package gitage

import (
	"os"
	"path/filepath"
)

func dir(path string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if path != "" {
		wd, err = filepath.Abs(path)
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(wd, ".gitage"), nil
}
