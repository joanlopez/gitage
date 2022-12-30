package main

import (
	"os"

	"github.com/joanlopez/gitage/cmd/gitage/bootstrap"
	"github.com/spf13/afero"
)

func main() {
	bootstrap.Run(afero.NewOsFs(), os.Stdout, os.Args[1:]...)
}
