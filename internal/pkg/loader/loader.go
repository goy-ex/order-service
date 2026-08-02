package loader

import (
	"errors"
	"fmt"
	"io"
)

type Fetcher interface {
	Fetch() (io.ReadCloser, error)
}

type FetcherFunc func() (io.ReadCloser, error)

func (f FetcherFunc) Fetch() (io.ReadCloser, error) {
	return f()
}

type Decoder[T any] interface {
	Decode(r io.Reader) (T, error)
}

type DecoderFunc[T any] func(r io.Reader) (T, error)

func (f DecoderFunc[T]) Decode(r io.Reader) (T, error) {
	return f(r)
}

type Loader[T any] struct {
	fetcher Fetcher
	decoder Decoder[T]
}

func NewLoader[T any](fetcher Fetcher, decoder Decoder[T]) *Loader[T] {
	return &Loader[T]{
		fetcher: fetcher,
		decoder: decoder,
	}
}

// Load provides.
//
//nolint:nonamedreturns // named return required to join the deferred error
func (l *Loader[T]) Load() (res T, err error) {
	fetched, err := l.fetcher.Fetch()
	if err != nil {
		err = fmt.Errorf("failed to fetch: %w", err)

		return
	}

	defer func() {
		closeErr := fetched.Close()
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to close: %w", closeErr))
		}
	}()

	res, err = l.decoder.Decode(fetched)
	if err != nil {
		err = fmt.Errorf("failed to decode: %w", err)

		return
	}

	return
}
