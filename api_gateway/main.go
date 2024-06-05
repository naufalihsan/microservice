package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/naufalihsan/msvc-api-gateway/gateway"
	common "github.com/naufalihsan/msvc-common"
	"github.com/naufalihsan/msvc-common/discovery"
	"github.com/naufalihsan/msvc-common/discovery/consul"
)

var (
	httpAddress   = common.EnvString("HTTP_ADDR", ":8000")
	consulAddress = common.EnvString("CONSUL_ADDR", "localhost:8500")
)

func main() {
	registry, err := consul.NewRegistry(consulAddress, common.ApiGatewayService)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceId := discovery.GenerateInstanceId(common.ApiGatewayService)

	if err := registry.Register(ctx, instanceId, common.ApiGatewayService, httpAddress); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceId, common.ApiGatewayService); err != nil {
				log.Fatal("failed to health check")
			}

			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceId, common.ApiGatewayService)

	orderGateaway := gateway.NewGrpcGateway(registry)

	mux := http.NewServeMux()
	httpHandler := NewHttpHandler(orderGateaway)
	httpHandler.registerRoutes(mux)

	log.Printf("Start http server at port %s", httpAddress)

	if err := http.ListenAndServe(httpAddress, mux); err != nil {
		log.Fatalf("Failed to start http server %v", err)
	}
}
