package fs

import (
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/go-git/go-billy/v5"
)

// Mkdir docs (TODO)
// - path MUST be an absolute path.
func Mkdir(fs billy.Filesystem, path string) error {
	return fs.MkdirAll(path, 0o755)
}

// Create docs (TODO)
// - path MUST be an absolute path.
func Create(fs billy.Filesystem, path string, contents []byte) error {
	f, err := fs.Create(path)
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
func Append(fs billy.Filesystem, path string, contents []byte) error {
	f, err := fs.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
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
func Read(fs billy.Filesystem, path string) ([]byte, error) {
	f, err := fs.Open(path)
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(f)

	if closeErr := f.Close(); err == nil {
		err = closeErr
	}

	return bytes, err
}

// RemoveAll docs (TODO)
// - path MUST be an absolute path.
func RemoveAll(fs billy.Filesystem, path string) error {
	return fs.Remove(path)
}

func WriteFile(fs billy.Filesystem, filename string, data []byte, perm os.FileMode) error {
	f, err := fs.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}

	if closeErr := f.Close(); err == nil {
		err = closeErr
	}

	return err
}

func Walk(fs billy.Filesystem, root string, walkFn filepath.WalkFunc) error {
	info, err := fs.Lstat(root)
	if err != nil {
		return walkFn(root, nil, err)
	}
	return walk(fs, root, info, walkFn)
}

func walk(fs billy.Filesystem, path string, info os.FileInfo, walkFn filepath.WalkFunc) error {
	err := walkFn(path, info, nil)
	if err != nil {
		if info.IsDir() && err == filepath.SkipDir {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return nil
	}

	names, err := readDirNames(fs, path)
	if err != nil {
		return walkFn(path, info, err)
	}

	for _, name := range names {
		filename := filepath.Join(path, name)
		fileInfo, err := fs.Lstat(filename)
		if err != nil {
			if err := walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir {
				return err
			}
		} else {
			err = walk(fs, filename, fileInfo, walkFn)
			if err != nil {
				if !fileInfo.IsDir() || err != filepath.SkipDir {
					return err
				}
			}
		}
	}
	return nil
}

func readDirNames(fs billy.Filesystem, dirname string) ([]string, error) {
	entries, err := fs.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}

	sort.Strings(names)
	return names, nil
}
