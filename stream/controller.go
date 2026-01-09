package stream

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"streammy/pipeline"
	"streammy/sources"
	"streammy/transforms"
)

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func StreamController() http.Handler {
	router := chi.NewRouter()

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// ---- QUERY PARAMS ----
		page := parseInt(r.URL.Query().Get("page"), 1)
		limit := parseInt(r.URL.Query().Get("limit"), 10)

		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 100 {
			limit = 10
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// ---- PIPELINE ----
		source := sources.HTTPSource(
			ctx,
			client,
			1, // one request per page
			limit,
			1, // no fan-out needed
			func(skip, take int) string {
				return "https://jsonplaceholder.typicode.com/posts" +
					"?_page=" + strconv.Itoa(page) +
					"&_limit=" + strconv.Itoa(take)
			},
		)

		decoded := transforms.DecodeJSONArray[Post]()(ctx, source)

		processed := pipeline.WorkerPool[Post, Post](
			4,
			func(ctx context.Context, p Post) (Post, error) {
				p.Title = "[STREAMED] " + p.Title
				return p, nil
			},
		)(ctx, decoded)

		_ = JSONStreamSink(ctx, w, processed)
	})

	return router
}

func parseInt(v string, def int) int {
	if v == "" {
		return def
	}
	if i, err := strconv.Atoi(v); err == nil {
		return i
	}
	return def
}
