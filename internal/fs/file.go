package fs

import (
	"io"
	"os"
)

type File interface {
	io.Closer
	io.Reader
	io.Writer

	Readdirnames(n int) ([]string, error)

	Stat() (os.FileInfo, error)
}
