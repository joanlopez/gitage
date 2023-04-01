package main

import (
	"os"

	"github.com/go-git/go-billy/v5/osfs"

	"github.com/joanlopez/gitage/cmd/gitage/bootstrap"
	"github.com/joanlopez/gitage/internal/log"
)

func main() {
	ctx := log.Ctx(os.Stdout)
	bootstrap.Run(ctx, osfs.New(""), os.Args[1:]...)
}
