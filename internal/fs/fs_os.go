package fs

import (
	"os"
)

var _ FS = &OsFS{}

type OsFS struct{}

func NewOsFs() FS {
	return &OsFS{}
}

func (OsFS) Create(name string) (File, error) {
	f, e := os.Create(name)
	if f == nil {
		// while this looks strange, we need to return a bare nil (of type nil) not
		// a nil value of type *os.File or nil won't be nil
		return nil, e
	}
	return f, e
}

func (OsFS) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (OsFS) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (OsFS) Open(name string) (File, error) {
	f, e := os.Open(name)
	if f == nil {
		// while this looks strange, we need to return a bare nil (of type nil) not
		// a nil value of type *os.File or nil won't be nil
		return nil, e
	}
	return f, e
}

func (OsFS) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	f, e := os.OpenFile(name, flag, perm)
	if f == nil {
		// while this looks strange, we need to return a bare nil (of type nil) not
		// a nil value of type *os.File or nil won't be nil
		return nil, e
	}
	return f, e
}

func (OsFS) Remove(name string) error {
	return os.Remove(name)
}

func (OsFS) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (OsFS) Rename(oldname, newname string) error {
	return os.Rename(oldname, newname)
}

func (OsFS) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(name)
}

func (OsFS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
