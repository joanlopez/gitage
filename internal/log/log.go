package log

import (
	"context"
	"fmt"
	"io"
	"os"
)

type writerKey struct{}

func Ctx(w io.Writer) context.Context {
	return context.WithValue(context.Background(), writerKey{}, Writer{w})
}

func For(ctx context.Context) Writer {
	if w, ok := ctx.Value(writerKey{}).(Writer); ok {
		return w
	}

	return Writer{os.Stdout}
}

type Writer struct {
	io.Writer
}

func (w Writer) Print(a ...any) {
	_, _ = fmt.Fprint(w.Writer, a...)
}

func (w Writer) Printf(format string, a ...any) {
	_, _ = fmt.Fprintf(w.Writer, format, a...)
}

func (w Writer) Println(a ...any) {
	_, _ = fmt.Fprintln(w.Writer, a...)
}
