package sources

import (
	"context"
	"net/http"
	"sync"
)

func HTTPSource(
	ctx context.Context,
	client *http.Client,
	pages int,
	pageSize int,
	concurrency int,
	buildURL func(skip, take int) string,
) <-chan *http.Response {

	out := make(chan *http.Response)

	go func() {
		defer close(out)

		var wg sync.WaitGroup
		sem := make(chan struct{}, concurrency)

		for i := 0; i < pages; i++ {
			skip := i * pageSize
			sem <- struct{}{}
			wg.Add(1)

			go func(skip int) {
				defer wg.Done()
				defer func() { <-sem }()

				req, err := http.NewRequestWithContext(
					ctx,
					http.MethodGet,
					buildURL(skip, pageSize),
					nil,
				)
				if err != nil {
					return
				}

				resp, err := client.Do(req)
				if err != nil {
					return
				}

				select {
				case out <- resp:
				case <-ctx.Done():
					resp.Body.Close()
				}
			}(skip)
		}

		wg.Wait()
	}()

	return out
}
