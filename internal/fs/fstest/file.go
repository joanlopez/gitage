package fstest

import (
	"golang.org/x/tools/txtar"
)

type File struct {
	txtar.File
}

func NewFile(name string, data []byte) *File {
	return &File{
		File: txtar.File{
			Name: name,
			Data: data,
		},
	}
}
