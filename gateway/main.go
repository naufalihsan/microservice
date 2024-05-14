package main

import (
	"log"
	"net/http"
)

const (
	httpAddress = ":8080"
)

func main() {
	mux := http.NewServeMux()

	httpHandler := NewHttpHandler()
	httpHandler.registerRoutes(mux)

	if err := http.ListenAndServe(httpAddress, mux); err != nil {
		log.Fatal("Failed to start http server")
	}
}
