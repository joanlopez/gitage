package fs

import (
	stdfs "io/fs"
	"os"
	"strings"

	"github.com/spf13/afero"
	"golang.org/x/tools/txtar"
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

	if _, err := f.Write(contents); err != nil {
		return err
	}

	return f.Close()
}

func rootify(path string) string {
	if strings.HasPrefix(path, root) || strings.HasPrefix(path, cwd) {
		return path
	}
	return root + path
}

func FromTxtar(archive *txtar.Archive) (FS, error) {
	fs := afero.NewMemMapFs()
	for _, f := range archive.Files {
		if err := afero.WriteFile(fs, f.Name, f.Data, 0644); err != nil {
			return nil, err
		}
	}

	return fs, nil
}

func ToTxtar(fs FS) (*txtar.Archive, error) {
	a := new(txtar.Archive)
	err := afero.Walk(fs, "/", func(path string, info stdfs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		contents, err := afero.ReadFile(fs, path)
		if err != nil {
			return err
		}

		a.Files = append(a.Files, txtar.File{
			Name: path,
			Data: contents,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return a, nil
}
