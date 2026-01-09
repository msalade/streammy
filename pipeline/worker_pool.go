package pipeline

import "context"

func WorkerPool[A any, B any](
	workers int,
	fn func(context.Context, A) (B, error),
) Stage[A, B] {
	return func(ctx context.Context, in <-chan A) <-chan B {
		out := make(chan B)

		go func() {
			defer close(out)

			sem := make(chan struct{}, workers)
			for item := range in {
				sem <- struct{}{}
				go func(item A) {
					defer func() { <-sem }()
					result, err := fn(ctx, item)
					if err != nil {
						return
					}

					select {
					case out <- result:
					case <-ctx.Done():
					}
				}(item)
			}

			for range workers {
				sem <- struct{}{}
			}
		}()

		return out
	}
}
