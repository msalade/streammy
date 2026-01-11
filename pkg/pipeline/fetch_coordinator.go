package pipeline

import "sync"

type Fetch[T any] func(skip, take int) ([]T, error)

func FetchCoordinator[T any](
	total, take, parallelism int,
	fetch Fetch[T],
) <-chan T {
	tasks := make(chan int)
	out := make(chan T, take*parallelism)

	go func() {
		for skip := 0; skip < total; skip += take {
			tasks <- skip
		}
		close(tasks)
	}()

	var wg sync.WaitGroup
	wg.Add(parallelism)

	for range parallelism {
		go func() {
			defer wg.Done()
			for skip := range tasks {
				items, err := fetch(skip, take)
				if err != nil {
					continue
				}
				for _, item := range items {
					out <- item
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
