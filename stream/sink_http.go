package stream

import (
	"context"
	"encoding/json"
	"io"
)

func JSONStreamSink[T any](ctx context.Context, w io.Writer, in <-chan T) error {
	enc := json.NewEncoder(w)

	if _, err := w.Write([]byte("[")); err != nil {
		return err
	}

	first := true
	for item := range in {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if !first {
			if _, err := w.Write([]byte(",")); err != nil {
				return err
			}
		}
		first = false

		if err := enc.Encode(item); err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("]"))
	return err
}
