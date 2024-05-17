package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/naufalihsan/msvc-common"
	pb "github.com/naufalihsan/msvc-common/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	httpAddress          = common.EnvString("HTTP_ADDR", ":8000")
	ordersServiceAddress = "localhost:3000" // temp port
)

func main() {
	conn, err := grpc.NewClient(ordersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	log.Printf("Dial orders service gRPC server at port %s", ordersServiceAddress)

	client := pb.NewOrderServiceClient(conn)

	mux := http.NewServeMux()
	httpHandler := NewHttpHandler(client)
	httpHandler.registerRoutes(mux)

	log.Printf("Start http server at port %s", httpAddress)

	if err := http.ListenAndServe(httpAddress, mux); err != nil {
		log.Fatal("Failed to start http server %v", err)
	}
}
