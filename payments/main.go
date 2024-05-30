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
	grpcAddress   = common.EnvString("GRPC_ADDR", "localhost:3001")
	consulAddress = common.EnvString("CONSUL_ADDR", "localhost:8500")
	amqpUser      = common.EnvString("AMQP_USER", "guest")
	amqpPass      = common.EnvString("AMQP_PASS", "guest")
	amqpHost      = common.EnvString("AMQP_HOST", "localhost")
	amqpPort      = common.EnvString("AMQP_PORT", "5672")
)

func main() {
	registry, err := consul.NewRegistry(consulAddress, common.PaymentService)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceId := discovery.GenerateInstanceId(common.PaymentService)

	if err := registry.Register(ctx, instanceId, common.PaymentService, grpcAddress); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceId, common.PaymentService); err != nil {
				log.Fatal("failed to health check")
			}

			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceId, common.PaymentService)

	channel, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		channel.Close()
	}()

	service := NewService()

	consumer := NewConsumer(service)
	go consumer.Listen(channel)

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	log.Printf("Start gRPC server at port %s", grpcAddress)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err.Error())
	}
}