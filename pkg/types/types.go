package types

import "context"

type Source[T any] interface {
	Stream(ctx context.Context) (<-chan T, <-chan error)
}

type Transformer[T any, U any] func(ctx context.Context, in <-chan T) <-chan U

type Sink[T any] interface {
	Write(ctx context.Context, in <-chan T) error
}
