package main

import (
	"log"
	"net/http"

	"streammy/stream"
)

func main() {
	r := stream.StreamController()

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
