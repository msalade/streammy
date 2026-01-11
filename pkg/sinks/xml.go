package sinks

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
)

type XML[T any] struct {
	Writer  io.Writer
	RootTag string
	ItemTag string
	Indent  bool
}

func (s *XML[T]) Write(ctx context.Context, in <-chan T) error {
	flusher, _ := s.Writer.(interface{ Flush() })

	if s.RootTag == "" {
		s.RootTag = "root"
	}
	if s.ItemTag == "" {
		s.ItemTag = "item"
	}
	if _, err := fmt.Fprintf(s.Writer, "<%s>\n", s.RootTag); err != nil {
		return err
	}

	enc := xml.NewEncoder(s.Writer)
	if s.Indent {
		enc.Indent("", "  ")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case item, ok := <-in:
			if !ok {
				if _, err := fmt.Fprintf(s.Writer, "</%s>\n", s.RootTag); err != nil {
					return err
				}
				return nil
			}

			wrapped := struct {
				XMLName xml.Name
				Value   T `xml:",inline"`
			}{
				XMLName: xml.Name{Local: s.ItemTag},
				Value:   item,
			}

			if err := enc.Encode(wrapped); err != nil {
				return err
			}

			if flusher != nil {
				flusher.Flush()
			}
		}
	}
}
