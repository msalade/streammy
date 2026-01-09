package pipeline

import "context"

type Ordered[T any] struct {
	Index int
	Value T
}

func PreserveOrder[T any]() Stage[Ordered[T], T] {
	return func(ctx context.Context, in <-chan Ordered[T]) <-chan T {
		out := make(chan T)

		go func() {
			defer close(out)

			buffer := map[int]T{}
			next := 0

			for item := range in {
				buffer[item.Index] = item.Value

				for {
					v, ok := buffer[next]
					if !ok {
						break
					}
					delete(buffer, next)
					out <- v
					next++
				}
			}
		}()

		return out
	}
}
