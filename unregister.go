package gitage

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5"

	"github.com/joanlopez/gitage/internal/fs"
	"github.com/joanlopez/gitage/internal/log"
)

// Unregister docs (TODO)
// - path MUST be an absolute path.
func Unregister(ctx context.Context, f billy.Filesystem, path string, recipients ...string) error {
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

	log.For(ctx).Println("Unregistering recipients...")

	contents, err := fs.Read(f, recipientsFilepath)
	if err != nil {
		return err
	}

	var result []byte

	fileScanner := bufio.NewScanner(bytes.NewBuffer(contents))
	for fileScanner.Scan() {
		next := fileScanner.Text()
		if !present(next, recipients...) {
			result = append(result, []byte(next)...)
			result = append(result, []byte("\n")...)
		}
	}

	if err := fs.Create(f, recipientsFilepath, result); err != nil {
		return err
	}

	log.For(ctx).Println("Recipients unregistered with success!")

	return nil
}

func present(r string, recipients ...string) bool {
	for _, recipient := range recipients {
		if r == recipient {
			return true
		}
	}

	return false
}
