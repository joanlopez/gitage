package archive

import (
	"golang.org/x/tools/txtar"
)

type Archive struct {
	ref   map[string]*txtar.File
	inner *txtar.Archive
}

func Empty() *Archive {
	a := new(Archive)
	a.init()
	return a
}

func Parse(data []byte) *Archive {
	inner := txtar.Parse(data)
	ref := make(map[string]*txtar.File, len(inner.Files))
	for i := range inner.Files {
		ref[inner.Files[i].Name] = &inner.Files[i]
	}
	return &Archive{
		ref:   ref,
		inner: inner,
	}
}

func (a *Archive) Add(f *File) {
	a.init()
	a.inner.Files = append(a.inner.Files, f.File)
	a.ref[f.Name] = &a.inner.Files[len(a.inner.Files)-1]
}

func (a *Archive) Get(path string) *File {
	a.init()
	if f, ok := a.ref[path]; ok {
		return &File{File: *f}
	}
	return nil
}

func (a *Archive) Exists(path string) bool {
	a.init()
	_, ok := a.ref[path]
	return ok
}

func (a *Archive) Format() []byte {
	a.init()
	return txtar.Format(a.inner)
}

func (a *Archive) Files() <-chan *File {
	a.init()
	ch := make(chan *File)
	go func() {
		for _, f := range a.inner.Files {
			ch <- &File{File: f}
		}
		close(ch)
	}()
	return ch
}

func (a *Archive) init() {
	if a.ref == nil {
		a.ref = make(map[string]*txtar.File)
	}

	if a.inner == nil {
		a.inner = &txtar.Archive{}
	}
}
