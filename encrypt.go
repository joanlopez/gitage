package gitage

import (
	"bytes"
	"context"
	"io"
	stdfs "io/fs"
	"path/filepath"

	"filippo.io/age"
	"github.com/spf13/afero"

	"github.com/joanlopez/gitage/internal/fs"
)

// EncryptAll encrypts all files in the specified path,
// so it is equivalent to calling EncryptFile for each
// file in the given path, recursively.
//
// It skips directories (files are encrypted individually)
// and encrypted files (files with the .age extension) to
// avoid double encryption.
func EncryptAll(ctx context.Context, f fs.FS, path string, recipients ...age.Recipient) error {
	return afero.Walk(f, path, func(path string, info stdfs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip encrypted files
		if filepath.Ext(path) == Ext {
			return nil
		}

		return EncryptFile(ctx, f, path, recipients...)
	})
}

// EncryptFile encrypts the file present at the given
// path, within the given file-system, using the given
// recipients.
//
// In comparison to Encrypt, it replaces the plain file
// with the encrypted one (with the .age extension).
//
// So, assuming it can be called with a non-transactional
// file-system, use it with care. An unsuccessful operation
// will leave the file-system in an inconsistent state.
func EncryptFile(ctx context.Context, f fs.FS, path string, recipients ...age.Recipient) error {
	file, err := f.Open(path)
	if err != nil {
		return err
	}

	read, err := afero.ReadAll(file)
	if err != nil {
		return err
	}
	file.Close()

	if err = fs.RemoveAll(f, path); err != nil {
		return err
	}

	toWrite, err := Encrypt(ctx, read, recipients...)
	if err != nil {
		return err
	}

	agedPath := path + Ext

	return fs.Create(f, agedPath, toWrite)
}

// Encrypt encrypts the given plaintext using the given
// recipients and 'age' encryption tool (Go library).
func Encrypt(_ context.Context, plaintext []byte, recipients ...age.Recipient) ([]byte, error) {
	buff := new(bytes.Buffer)

	w, err := age.Encrypt(buff, recipients...)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(w, bytes.NewReader(plaintext)); err != nil {
		return nil, err
	}

	if err = w.Close(); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
