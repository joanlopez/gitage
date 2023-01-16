package fs

import (
	"fmt"
	stdfs "io/fs"
	"os"

	"github.com/spf13/afero"

	"github.com/joanlopez/gitage/internal/fs/archive"
)

type FS afero.Fs

// Mkdir docs (TODO)
// - path MUST be an absolute path.
func Mkdir(fs FS, path string) error {
	return fs.MkdirAll(normalize(path), 0o755)
}

// Create docs (TODO)
// - path MUST be an absolute path.
func Create(fs FS, path string, contents []byte) error {
	f, err := fs.Create(normalize(path))
	if err != nil {
		return err
	}

	if _, err := f.Write(contents); err != nil {
		return err
	}

	return f.Close()
}

// Append docs (TODO)
// - path MUST be an absolute path.
func Append(fs FS, path string, contents []byte) error {
	f, err := fs.OpenFile(normalize(path), os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	if _, err := f.Write(contents); err != nil {
		return err
	}

	return f.Close()
}

// Read docs (TODO)
// - path MUST be an absolute path.
func Read(fs FS, path string) ([]byte, error) {
	f, err := fs.Open(normalize(path))
	if err != nil {
		return nil, err
	}

	contents, err := afero.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return contents, f.Close()
}

// RemoveAll docs (TODO)
// - path MUST be an absolute path.
func RemoveAll(fs FS, path string) error {
	return fs.RemoveAll(normalize(path))
}

// For more context, look at these links:
// - https://github.com/golang/go/issues/21782
// - https://github.com/spf13/afero/pull/302/files
func normalize(path string) string {
	return normalizeLongPath(path)
}

func FromArchive(a *archive.Archive) (FS, error) {
	fs := afero.NewMemMapFs()
	for f := range a.Files() {
		if err := afero.WriteFile(fs, f.Name, f.Data, 0o644); err != nil {
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

		// Skip directories,
		// cause Archive only represents files.
		if info.IsDir() {
			return nil
		}

		contents, err := afero.ReadFile(fs, normalize(path))
		if err != nil {
			return err
		}

		a.Add(archive.NewFile(path, contents))
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot build archive for fs: %w", err)
	}

	return a, nil
}
