package transforms

import (
	"context"
	"encoding/json"
	"net/http"
)

func DecodeJSONArray[T any]() func(context.Context, <-chan *http.Response) <-chan T {
	return func(ctx context.Context, in <-chan *http.Response) <-chan T {
		out := make(chan T)

		go func() {
			defer close(out)

			for resp := range in {
				dec := json.NewDecoder(resp.Body)

				tok, err := dec.Token()
				if err != nil || tok != json.Delim('[') {
					resp.Body.Close()
					continue
				}

				for dec.More() {
					var item T
					if err := dec.Decode(&item); err != nil {
						break
					}

					select {
					case out <- item:
					case <-ctx.Done():
						resp.Body.Close()
						return
					}
				}
				resp.Body.Close()
			}
		}()

		return out
	}
}
