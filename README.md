# Gitage [![GoDoc](https://godoc.org/gopkg.in/joanlopez/gitage?status.svg)](https://pkg.go.dev/github.com/joanlopez/gitage) [![Test](https://github.com/joanlopez/gitage/workflows/Test/badge.svg)](https://github.com/joanlopez/gitage/actions?query=workflow%3ATest)

**[Git](https://git-scm.com/)+[age](https://github.com/FiloSottile/age) = Gitage;** simple, modern and secure Git encryption
tool

Gitage is a CLI tool that can be used as a wrapper of Git CLI. 

It uses [`age`](https://github.com/FiloSottile/age) encryption tool to encrypt files before committing them to the repository.

### Credits

This project relies on [age](https://github.com/FiloSottile/age), thanks to [Filippo Valsorda](https://github.com/FiloSottile).

It also relies on [go-git](https://github.com/go-git/go-git) and its file-system abstraction [billy](https://github.com/go-git/go-billy),
so thanks to their contributors as well.

