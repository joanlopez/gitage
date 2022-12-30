package main

import (
	"github.com/joanlopez/gitage/internal/log"
	"os"

	"github.com/joanlopez/gitage/cmd/gitage/bootstrap"
	"github.com/spf13/afero"
)

func main() {
	ctx := log.Ctx(os.Stdout)
	bootstrap.Run(ctx, afero.NewOsFs(), os.Args[1:]...)
}
