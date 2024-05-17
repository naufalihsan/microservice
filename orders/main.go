package main

import (
	"context"
	"log"
	"net"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/naufalihsan/msvc-common"
	"google.golang.org/grpc"
)

var (
	grpcAddress = common.EnvString("GRPC_ADDR", "localhost:3000")
)

func main() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	NewGrpcHandler(grpcServer)

	store := NewStore()
	svc := NewService(store)
	svc.CreateOrder(context.Background())

	log.Printf("Start gRPC server at port %s", grpcAddress)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err.Error())
	}
}
