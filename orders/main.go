package main

import (
	"context"
	"log"
	"net"
	"time"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/naufalihsan/msvc-common"
	"github.com/naufalihsan/msvc-common/broker"
	"github.com/naufalihsan/msvc-common/discovery"
	"github.com/naufalihsan/msvc-common/discovery/consul"
	"google.golang.org/grpc"
)

var (
	grpcAddress   = common.EnvString("GRPC_ADDR", "localhost:3000")
	jaegerAddress = common.EnvString("JAEGER_ADDR", "localhost:4318")
	consulAddress = common.EnvString("CONSUL_ADDR", "localhost:8500")
	amqpUser      = common.EnvString("AMQP_USER", "guest")
	amqpPass      = common.EnvString("AMQP_PASS", "guest")
	amqpHost      = common.EnvString("AMQP_HOST", "localhost")
	amqpPort      = common.EnvString("AMQP_PORT", "5672")
)

func main() {
	if err := common.SetGlobalTracer(context.TODO(), common.OrderService, jaegerAddress); err != nil {
		log.Fatal(err)
	}

	registry, err := consul.NewRegistry(consulAddress, common.OrderService)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceId := discovery.GenerateInstanceId(common.OrderService)

	if err := registry.Register(ctx, instanceId, common.OrderService, grpcAddress); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceId, common.OrderService); err != nil {
				log.Fatal("failed to health check")
			}

			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceId, common.OrderService)

	channel, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		channel.Close()
	}()

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	store := NewStore()
	service := NewService(store)
	serviceMiddleware := NewTelemetry(service)

	consumer := NewConsumer(serviceMiddleware)
	go consumer.Listen(channel, instanceId)

	NewGrpcHandler(grpcServer, serviceMiddleware, channel)
	log.Printf("Start gRPC server at port %s", grpcAddress)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err.Error())
	}
}
