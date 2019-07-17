package ctxful

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReaderOK(t *testing.T) {
	r := NewReader(context.Background(), readerFunc(func(b []byte) (int, error) {
		copy(b, []byte("foo"))
		return len("foo"), nil
	}))
	b := make([]byte, 999)
	n, err := r.Read(b)
	assert.NoError(t, err)
	assert.Equal(t, n, len("foo"))
	assert.Equal(t, "foo", string(b[:n]))
}

func TestReaderCancel(t *testing.T) {
	read := make(chan struct{})
	returned := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	r := NewReader(ctx, readerFunc(func(b []byte) (int, error) {
		defer close(read)
		<-returned

		copy(b, []byte("foo"))
		return len("foo"), nil
	}))

	b := make([]byte, 999)
	_, err := r.Read(b)
	assert.Error(t, err)

	close(returned)
	<-read
}

func TestWriterOK(t *testing.T) {
	w := NewWriter(context.Background(), writerFunc(func(b []byte) (int, error) {
		return len(b), nil
	}))
	b := []byte("foo")
	n, err := w.Write(b)
	assert.NoError(t, err)
	assert.Equal(t, n, len("foo"))
	assert.Equal(t, "foo", string(b[:n]))
}

func TestWriterCancel(t *testing.T) {
	wrote := make(chan struct{})
	returned := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	w := NewWriter(ctx, writerFunc(func(b []byte) (int, error) {
		defer close(wrote)
		<-returned

		return len(b), nil
	}))

	b := []byte("foo")
	_, err := w.Write(b)
	assert.Error(t, err)

	close(returned)
	<-wrote
}
