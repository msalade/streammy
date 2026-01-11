package sinks

import (
	"context"
	"encoding/json"
	"io"
)

type JSONL[T any] struct {
	Writer io.Writer
}

func (s *JSONL[T]) Write(ctx context.Context, in <-chan T) error {
	enc := json.NewEncoder(s.Writer)
	enc.SetEscapeHTML(false)
	flusher, _ := s.Writer.(interface{ Flush() })

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case item, ok := <-in:
			if !ok {
				return nil
			}
			if err := enc.Encode(item); err != nil {
				return err
			}
			if flusher != nil {
				flusher.Flush()
			}
		}
	}
}
