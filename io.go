// Package ctxful passes contexts to things that don't take contexts.
package ctxful

import (
	"context"
	"io"
)

func NewReader(ctx context.Context, r io.Reader) io.Reader {
	return readerFunc(readWriteContextFunc(ctx, r.Read))
}

func NewWriter(ctx context.Context, w io.Writer) io.Writer {
	return writerFunc(readWriteContextFunc(ctx, w.Write))
}

func readWriteContextFunc(ctx context.Context, op func([]byte) (int, error)) func([]byte) (int, error) {
	return func(b []byte) (int, error) {
		got := make(chan struct{})

		var n int
		var err error

		go func() {
			n, err = op(b)
			close(got)
		}()

		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-got:
			return n, err
		}
	}
}

type readerFunc func([]byte) (n int, err error)

func (f readerFunc) Read(b []byte) (n int, err error) {
	return f(b)
}

type writerFunc func([]byte) (n int, err error)

func (f writerFunc) Write(b []byte) (n int, err error) {
	return f(b)
}
