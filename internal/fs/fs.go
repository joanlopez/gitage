package fs

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sort"
)

type Fs interface {
	// Create creates a file in the filesystem, returning the file and an error, if any happens.
	Create(name string) (File, error)

	// Mkdir creates a directory in the filesystem, return an error if any happens.
	Mkdir(name string, perm os.FileMode) error

	// MkdirAll creates a directory path and all parents that does not exist yet.
	MkdirAll(path string, perm os.FileMode) error

	// Open opens a file, returning it or an error, if any happens.
	Open(name string) (File, error)

	// OpenFile opens a file using the given flags and the given mode.
	OpenFile(name string, flag int, perm os.FileMode) (File, error)

	// Remove removes a file identified by name, returning an error, if any happens.
	Remove(name string) error

	// RemoveAll removes a directory path and any children it contains.
	// It does not fail if the path does not exist (return nil).
	RemoveAll(path string) error

	// Stat returns a FileInfo describing the named file, or an error, if any happens.
	Stat(name string) (os.FileInfo, error)

	// Lstat returns a FileInfo describing the named file, or an error, if any happens.
	Lstat(name string) (os.FileInfo, error)
}

// Mkdir docs (TODO)
// - path MUST be an absolute path.
func Mkdir(fs Fs, path string) error {
	return fs.MkdirAll(normalize(path), 0o755)
}

// Create docs (TODO)
// - path MUST be an absolute path.
func Create(fs Fs, path string, contents []byte) error {
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
func Append(fs Fs, path string, contents []byte) error {
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
func Read(fs Fs, path string) ([]byte, error) {
	f, err := fs.Open(normalize(path))
	if err != nil {
		return nil, err
	}

	contents, err := ReadAll(f)
	if err != nil {
		return nil, err
	}

	return contents, f.Close()
}

func ReadAll(r io.Reader) ([]byte, error) {
	return readAll(r, bytes.MinRead)
}

func ReadFile(fs Fs, filename string) ([]byte, error) {
	f, err := fs.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// It's a good but not certain bet that FileInfo will tell us exactly how much to
	// read, so let's try it but be prepared for the answer to be wrong.
	var n int64

	if fi, err := f.Stat(); err == nil {
		// Don't preallocate a huge buffer, just in case.
		if size := fi.Size(); size < 1e9 {
			n = size
		}
	}
	// As initial capacity for readAll, use n + a little extra in case Size is zero,
	// and to avoid another allocation after Read has filled the buffer.  The readAll
	// call will read into its allocated internal buffer cheaply.  If the size was
	// wrong, we'll either waste some space off the end or reallocate as needed, but
	// in the overwhelmingly common case we'll get it just right.
	return readAll(f, n+bytes.MinRead)
}

// RemoveAll docs (TODO)
// - path MUST be an absolute path.
func RemoveAll(fs Fs, path string) error {
	return fs.RemoveAll(normalize(path))
}

func WriteFile(fs Fs, filename string, data []byte, perm os.FileMode) error {
	f, err := fs.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// For more context, look at these links:
// - https://github.com/golang/go/issues/21782
// - https://github.com/spf13/afero/pull/302/files
func normalize(path string) string {
	const longPath = 180
	if len(path) < longPath {
		return path
	}
	return normalizeLongPath(path)
}

func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, capacity))
	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}

func Walk(fs Fs, root string, walkFn filepath.WalkFunc) error {
	info, err := fs.Lstat(root)
	if err != nil {
		return walkFn(root, nil, err)
	}
	return walk(fs, root, info, walkFn)
}

func walk(fs Fs, path string, info os.FileInfo, walkFn filepath.WalkFunc) error {
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

func readDirNames(fs Fs, dirname string) ([]string, error) {
	f, err := fs.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}
