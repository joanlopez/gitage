package fs

import (
	stdfs "io/fs"
	"os"
	"strings"

	"github.com/spf13/afero"

	"github.com/joanlopez/gitage/internal/fs/archive"
)

const root = "/"
const cwd = "./"

type FS afero.Fs

func Mkdir(fs FS, path string) error {
	return fs.MkdirAll(rootify(path), 0755)
}

func Create(fs FS, path string, contents []byte) error {
	f, err := fs.Create(rootify(path))
	if err != nil {
		return err
	}

	if _, err := f.Write(contents); err != nil {
		return err
	}

	return f.Close()
}

func Append(fs FS, path string, contents []byte) error {
	f, err := fs.OpenFile(rootify(path), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if _, err := f.Write(contents); err != nil {
		return err
	}

	return f.Close()
}

func Read(fs FS, path string) ([]byte, error) {
	f, err := fs.Open(rootify(path))
	if err != nil {
		return nil, err
	}

	contents, err := afero.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return contents, f.Close()
}

func RemoveAll(fs FS, path string) error {
	return fs.RemoveAll(normalize(path))
}

// For more context, look at these links:
// - https://github.com/golang/go/issues/21782
// - https://github.com/spf13/afero/pull/302/files
func normalize(path string) string {
	return normalizeLongPath(path)
}

func rootify(path string) string {
	if strings.HasPrefix(path, root) || strings.HasPrefix(path, cwd) {
		return path
	}
	return root + path
}

func FromArchive(a *archive.Archive) (FS, error) {
	fs := afero.NewMemMapFs()
	for f := range a.Files() {
		if err := afero.WriteFile(fs, f.Name, f.Data, 0644); err != nil {
			return nil, err
		}
	}

	return fs, nil
}

func ToArchive(fs FS) (*archive.Archive, error) {
	a := archive.Empty()
	err := afero.Walk(fs, "/", func(path string, info stdfs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		contents, err := afero.ReadFile(fs, path)
		if err != nil {
			return err
		}

		a.Add(archive.NewFile(path, contents))
		return nil
	})

	if err != nil {
		return nil, err
	}

	return a, nil
}
