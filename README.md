# Gitage

**[Git](https://git-scm.com/)+[age](https://github.com/FiloSottile/age) = Gitage;** simple, modern and secure Git encryption
tool

Gitage is a CLI tool that can be used as a wrapper of Git CLI. 

It uses [`age`](https://github.com/FiloSottile/age) encryption tool to encrypt files before committing them to the repository.

### Credits

This project relies on [age](https://github.com/FiloSottile/age), thanks to [Filippo Valsorda](https://github.com/FiloSottile).

The file-system abstraction used in this project is a subset of [afero](https://github.com/spf13/afero), mainly focused
on "in-memory" and "os" implementations, and adapted for convenience, thanks to [Steve Francia](https://github.com/spf13).

