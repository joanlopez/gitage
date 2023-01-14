package main

import (
	"os"

	"github.com/spf13/afero"

	"github.com/joanlopez/gitage/cmd/gitage/bootstrap"
	"github.com/joanlopez/gitage/internal/log"
)

func main() {
	ctx := log.Ctx(os.Stdout)
	bootstrap.Run(ctx, afero.NewOsFs(), os.Args[1:]...)
}
