package pipeline

import "context"

type Source[T any] func(ctx context.Context) <-chan T
type Stage[A any, B any] func(ctx context.Context, in <-chan A) <-chan B
type Sink[T any] func(ctx context.Context, in <-chan T) error
