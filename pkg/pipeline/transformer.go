package pipeline

import "context"

func Map[T any, U any](f func(T) U) func(ctx context.Context, in <-chan T) <-chan U {
	return func(ctx context.Context, in <-chan T) <-chan U {
		out := make(chan U)
		go func() {
			defer close(out)
			for item := range in {
				select {
				case <-ctx.Done():
					return
				case out <- f(item):
				}
			}
		}()
		return out
	}
}

func Filter[T any](predicate func(T) bool) func(ctx context.Context, in <-chan T) <-chan T {
	return func(ctx context.Context, in <-chan T) <-chan T {
		out := make(chan T)
		go func() {
			defer close(out)
			for item := range in {
				if !predicate(item) {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case out <- item:
				}
			}
		}()
		return out
	}
}
