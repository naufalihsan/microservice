package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/naufalihsan/msvc-common"
)

var (
	httpAddress = common.EnvString("HTTP_ADDR", ":8080")
)

func main() {
	mux := http.NewServeMux()

	httpHandler := NewHttpHandler()
	httpHandler.registerRoutes(mux)

	log.Printf("Start http server at port %s", httpAddress)

	if err := http.ListenAndServe(httpAddress, mux); err != nil {
		log.Fatal("Failed to start http server")
	}
}
