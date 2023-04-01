package gitage

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"

	"github.com/joanlopez/gitage/internal/fs"
	"github.com/joanlopez/gitage/internal/log"
)

// Init docs (TODO)
// - path MUST be an absolute path.
func Init(ctx context.Context, f billy.Filesystem, path string, recipients ...string) error {
	gitageDir := dir(path)
	info, err := f.Stat(gitageDir)
	if err == nil {
		if !info.IsDir() {
			log.For(ctx).Printf(`%s already exists as a file...
Please, remove it and try again.`, gitageDir)
			return nil
		}

		if info.IsDir() {
			log.For(ctx).Printf("%s already exists...\nAre you in a Gitage repository already?\n", gitageDir)
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

	if err := initGitRepository(ctx, f, path); err != nil {
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

func initGitRepository(ctx context.Context, f billy.Filesystem, path string) error {
	log.For(ctx).Printf("Initializing Git repository at %s...\n", path)

	root, err := f.Chroot(path)
	if err != nil {
		return err
	}

	dot, err := root.Chroot(git.GitDirName)
	if err != nil {
		return err
	}

	s := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())

	r, err := git.Init(s, root)
	if err != nil {
		// We're done, with success!
		if errors.Is(err, git.ErrRepositoryAlreadyExists) {
			return nil
		}
		return err
	}

	// Trick to set default branch to main
	err = s.SetReference(plumbing.NewSymbolicReference(plumbing.HEAD, "refs/heads/main"))
	if err != nil {
		return err
	}

	// And update the default branch in config
	cfg, err := r.Config()
	cfg.Init.DefaultBranch = "main"
	return s.SetConfig(cfg)
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
