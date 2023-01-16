package gitage

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/joanlopez/gitage/internal/fs"
	"github.com/joanlopez/gitage/internal/log"
)

func Register(ctx context.Context, f fs.FS, path string, recipients ...string) error {
	gitageDir := dir(path)
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

	contents, err := fs.Read(f, recipientsFilepath)
	if err != nil {
		return err
	}

	var bytes []byte

	for _, r := range recipients {
		if !strings.Contains(string(contents), r) {
			bytes = append(bytes, []byte(r)...)
			bytes = append(bytes, []byte("\n")...)
		}
	}

	if len(bytes) > 0 {
		if err := fs.Append(f, recipientsFilepath, recipientsBytes(recipients...)); err != nil {
			return err
		}
	}

	log.For(ctx).Println("Recipients registered with success!")

	return nil
}
