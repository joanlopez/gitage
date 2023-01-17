package fstest

import (
	"fmt"
	stdfs "io/fs"
	"os"

	"github.com/joanlopez/gitage/internal/fs"
)

const pathSepStr = string(os.PathSeparator)

func FsFromArchive(a *Archive) (fs.Fs, error) {
	memFs := fs.NewMemFs()
	for f := range a.Files() {
		if err := fs.WriteFile(memFs, f.Name, f.Data, 0o644); err != nil {
			return nil, err
		}
	}

	return memFs, nil
}

func FsToArchive(f fs.Fs) (*Archive, error) {
	a := EmptyArchive()
	err := fs.Walk(f, rootDir(), func(path string, info stdfs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories,
		// cause Archive only represents files.
		if info.IsDir() {
			return nil
		}

		contents, err := fs.ReadFile(f, path)
		if err != nil {
			return err
		}

		a.Add(NewFile(path, contents))
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot build archive for fs: %w", err)
	}

	return a, nil
}

func FsFromTxtarFile(pathToFile string) (fs.Fs, error) {
	b, err := os.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}

	f, err := FsFromArchive(ParseArchive(b))
	if err != nil {
		return nil, err
	}

	return f, nil
}
