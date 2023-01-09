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

// DecryptAll decrypts all files in the specified path,
// so it is equivalent to calling DecryptFile for each
// file in the given path, recursively.
//
// It skips directories (files are decrypted individually)
// and non-encrypted files (files without the .age extension)
// to avoid double decryption.
func DecryptAll(ctx context.Context, f fs.FS, path string, identities ...age.Identity) error {
	return afero.Walk(f, path, func(path string, info stdfs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip non-encrypted files
		if filepath.Ext(path) != Ext {
			return nil
		}

		err = DecryptFile(ctx, f, path, identities...)
		if err != nil {
			return err
		}

		return nil
	})
}

// DecryptFile decrypts the file present at the given
// path, within the given file-system, using the given
// identities.
//
// In comparison to Decrypt, it replaces the ciphered
// file with the decrypted one (w/out the .age extension).
//
// So, assuming it can be called with a non-transactional
// file-system, use it with care. An unsuccessful operation
// will leave the file-system in an inconsistent state.
func DecryptFile(ctx context.Context, f fs.FS, path string, identities ...age.Identity) error {
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

	toWrite, err := Decrypt(ctx, read, identities...)
	if err != nil {
		return err
	}

	path = path[:len(path)-len(Ext)]

	err = fs.Create(f, path, toWrite)
	if err != nil {
		return err
	}

	return nil
}

// Decrypt decrypts the given ciphertext using the given
// recipients and 'age' encryption tool (Go library).
func Decrypt(_ context.Context, ciphertext []byte, identities ...age.Identity) ([]byte, error) {
	buff := new(bytes.Buffer)

	r, err := age.Decrypt(bytes.NewReader(ciphertext), identities...)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(buff, r); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
