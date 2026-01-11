package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"streammy/pkg/pipeline"
	"streammy/pkg/sources"
	"streammy/pkg/types"

	"github.com/go-chi/chi/v5"
)

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func main() {
	r := chi.NewRouter()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	r.Get("/posts", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		total := 100
		take := 10
		parallelism := 4

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		var src types.Source[Post] = &pipeline.FetchSource[Post]{
			Fetch: func(skip, take int) ([]Post, error) {
				page := (skip / take) + 1
				url := "https://jsonplaceholder.typicode.com/posts?_page=" +
					strconv.Itoa(page) + "&_limit=" + strconv.Itoa(take)
				return sources.HTTPFetch[Post](ctx, client, url)
			},
			Total:       total,
			Take:        take,
			Parallelism: parallelism,
		}

		ch, errs := src.Stream(ctx)

		go func() {
			for err := range errs {
				if err != nil && err != context.Canceled {
					log.Println("source error:", err)
				}
			}
		}()

		sink := &pipeline.JSONSink[Post]{Writer: w}

		if err := sink.Write(ctx, ch); err != nil && err != context.Canceled {
			log.Println("sink error:", err)
		}
	})

	log.Println("HTTP server listening on http://localhost:8080/posts")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
