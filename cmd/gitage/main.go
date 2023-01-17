package main

import (
	"os"

	"github.com/joanlopez/gitage/cmd/gitage/bootstrap"
	"github.com/joanlopez/gitage/internal/fs"
	"github.com/joanlopez/gitage/internal/log"
)

func main() {
	ctx := log.Ctx(os.Stdout)
	bootstrap.Run(ctx, fs.NewOsFs(), os.Args[1:]...)
}
