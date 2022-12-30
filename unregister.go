package gitage

import (
	"context"
	"os"
	"path/filepath"

	"github.com/joanlopez/gitage/internal/fs"
	"github.com/joanlopez/gitage/internal/log"
)

func Unregister(ctx context.Context, f fs.FS, path string, recipients ...string) error {
	gitageDir, err := dir(path)
	if err != nil {
		return err
	}

	info, err := f.Stat(gitageDir)
	if (err != nil && os.IsNotExist(err)) || !info.IsDir() {
		log.For(ctx).Printf(`%s directory not found...
Are you in a Gitage repository?`, gitageDir)
		return nil
	}

	recipientsFilepath := filepath.Join(gitageDir, "recipients")

	info, err = f.Stat(recipientsFilepath)
	if (err != nil && os.IsNotExist(err)) || info.IsDir() {
		log.For(ctx).Printf(`%s file not found...
Are you in a Gitage repository?`, gitageDir)
		return nil
	}

	if err != nil {
		return err
	}

	log.For(ctx).Println("Registering recipients...")

	log.For(ctx).Println("Recipients registered with success!")

	return nil
}
