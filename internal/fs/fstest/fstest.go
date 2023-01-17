package fstest

import (
	"fmt"
	stdfs "io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joanlopez/gitage/internal/fs"
)

const pathSepStr = string(os.PathSeparator)

func FsFromArchive(a *Archive) (fs.Fs, error) {
	memFs := fs.NewMemFs()
	for f := range a.Files() {
		fName := Rootify(f.Name)
		if err := fs.WriteFile(memFs, fName, f.Data, 0o644); err != nil {
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

func FsFromPath(fsPath string) (fs.Fs, error) {
	memFs := fs.NewMemFs()

	err := filepath.Walk(fsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath := relPath(fsPath, path)

		if info.IsDir() {
			return fs.Mkdir(memFs, relPath)
		}

		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if err = fs.Create(memFs, relPath, contents); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return memFs, nil
}

func relPath(root, path string) string {
	path = strings.Replace(path, root, "", 1)
	if path == "" {
		path = "/"
	}

	return filepath.Clean(Rootify(path))
}

func FileContents(archive *Archive, f *File) []byte {
	file := archive.Get(f.Name)
	if file == nil {
		return nil
	}

	if runtime.GOOS == "windows" {
		return []byte(strings.ReplaceAll(string(file.Data), "\r\n", "\n"))
	}

	return file.Data
}
