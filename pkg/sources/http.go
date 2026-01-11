package sources

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

func HTTPFetch[T any](ctx context.Context, client *http.Client, url string) ([]T, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	tok, err := dec.Token()
	if err != nil || tok != json.Delim('[') {
		return nil, errors.New("expected JSON array")
	}

	var items []T
	for dec.More() {
		var item T
		if err := dec.Decode(&item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	_, err = dec.Token()
	return items, err
}
