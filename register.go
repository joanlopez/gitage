package gitage

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-billy/v5"

	"github.com/joanlopez/gitage/internal/fs"
	"github.com/joanlopez/gitage/internal/log"
)

// Register docs (TODO)
// - path MUST be an absolute path.
func Register(ctx context.Context, f billy.Filesystem, path string, recipients ...string) error {
	gitageDir := dir(path)
	info, err := f.Stat(gitageDir)
	if (err != nil && os.IsNotExist(err)) || !info.IsDir() {
		log.For(ctx).Printf("%s directory not found...\nAre you in a Gitage repository?\n", gitageDir)
		return nil
	}

	recipientsFilepath := filepath.Join(gitageDir, "recipients")

	info, err = f.Stat(recipientsFilepath)
	if (err != nil && os.IsNotExist(err)) || info.IsDir() {
		log.For(ctx).Printf("%s file not found...\nAre you in a Gitage repository?\n", gitageDir)
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
