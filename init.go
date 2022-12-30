package gitage

import (
	"context"
	"os"
	"path/filepath"

	"github.com/joanlopez/gitage/internal/fs"
	"github.com/joanlopez/gitage/internal/log"
)

func Init(ctx context.Context, f fs.FS, path string, recipients ...string) error {
	gitageDir, err := dir(path)
	if err != nil {
		return err
	}

	info, err := f.Stat(gitageDir)
	if err == nil {
		if !info.IsDir() {
			log.For(ctx).Printf(`%s already exists as a file...
Please, remove it and try again.`, gitageDir)
			return nil
		}

		if info.IsDir() {
			log.For(ctx).Printf(`%s already exists...
Are you in a Gitage repository already?`, gitageDir)
			return nil
		}
	}

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	log.For(ctx).Printf("Creating %s directory...\n", gitageDir)

	if err := fs.Mkdir(f, gitageDir); err != nil {
		return err
	}

	if err := fs.Create(f, filepath.Join(gitageDir, "config"), nil); err != nil {
		return err
	}

	if err := fs.Create(f, filepath.Join(gitageDir, "recipients"), recipientsBytes(recipients...)); err != nil {
		return err
	}

	log.For(ctx).Println("Gitage repository initialized with success!")

	return nil
}

func recipientsBytes(recipients ...string) []byte {
	if len(recipients) == 0 {
		return nil
	}

	var bytes []byte

	for _, r := range recipients {
		bytes = append(bytes, []byte(r)...)
		bytes = append(bytes, []byte("\n")...)
	}

	return bytes
}
