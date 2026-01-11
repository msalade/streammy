package pipeline

import (
	"context"
	"streammy/pkg/types"
)

type FetchSource[T any] struct {
	Fetch       Fetch[T]
	Total       int
	Take        int
	Parallelism int
}

func (f *FetchSource[T]) Stream(ctx context.Context) (<-chan T, <-chan error) {
	out := make(chan T)
	errs := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errs)

		fetchWithCtx := func(skip, take int) ([]T, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				return f.Fetch(skip, take)
			}
		}

		ch := FetchCoordinator(f.Total, f.Take, f.Parallelism, fetchWithCtx)
		for {
			select {
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			case item, ok := <-ch:
				if !ok {
					return
				}
				out <- item
			}
		}
	}()

	return out, errs
}

var _ types.Source[any] = (*FetchSource[any])(nil)
