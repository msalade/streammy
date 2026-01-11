package sinks

import (
	"context"
	"encoding/json"
	"io"
)

type JSONArray[T any] struct {
	Writer io.Writer
}

func (s *JSONArray[T]) Write(ctx context.Context, in <-chan T) error {
	enc := json.NewEncoder(s.Writer)
	enc.SetEscapeHTML(false)

	// open array
	if _, err := s.Writer.Write([]byte("[")); err != nil {
		return err
	}

	first := true
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case item, ok := <-in:
			if !ok {
				// close array
				_, err := s.Writer.Write([]byte("]"))
				return err
			}
			if !first {
				if _, err := s.Writer.Write([]byte(",")); err != nil {
					return err
				}
			}
			first = false
			if err := enc.Encode(item); err != nil {
				return err
			}
		}
	}
}
