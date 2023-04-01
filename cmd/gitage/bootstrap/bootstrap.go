package bootstrap

import (
	"context"

	"github.com/go-git/go-billy/v5"

	"github.com/joanlopez/gitage/cmd/gitage/bootstrap/cli"
	"github.com/joanlopez/gitage/internal/log"
)

func Run(ctx context.Context, fs billy.Filesystem, args ...string) {
	// Then we initialize a CLI with the given fs and out
	app := cli.New(ctx, fs)

	// Finally we run the CLI
	if err := app.Execute(args...); err != nil {
		log.For(ctx).Printf("Error: %s\n", err)
	}
}
