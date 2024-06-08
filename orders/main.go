package main

import (
	"context"
	"net"
	"time"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/naufalihsan/msvc-common"
	"github.com/naufalihsan/msvc-common/broker"
	"github.com/naufalihsan/msvc-common/discovery"
	"github.com/naufalihsan/msvc-common/discovery/consul"
	"github.com/naufalihsan/msvc-orders/gateway"
	"go.uber.org/zap"
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
	mongoUser     = common.EnvString("MONGO_USER", "root")
	mongoPass     = common.EnvString("MONGO_PASS", "example")
	mongoHost     = common.EnvString("MONGO_HOST", "localhost")
	mongoPort     = common.EnvString("MONGO_PORT", "27018")
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	if err := common.SetGlobalTracer(context.TODO(), common.OrderService, jaegerAddress); err != nil {
		logger.Fatal(err.Error())
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
				logger.Fatal("failed to health check", zap.Error(err))
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
		logger.Fatal(err.Error())
	}
	defer listener.Close()

	mongoClient, err := ConnectMongo(mongoUser, mongoPass, mongoHost, mongoPort)
	if err != nil {
		logger.Fatal("failed to connect to mongo db", zap.Error(err))
	}

	store := NewStore(mongoClient)
	inventoryGateaway := gateway.NewGrpcGateway(registry)
	service := NewService(store, inventoryGateaway)
	serviceMiddleware := NewLogging(NewTelemetry(service))

	consumer := NewConsumer(serviceMiddleware)
	go consumer.Listen(channel, instanceId)

	NewGrpcHandler(grpcServer, serviceMiddleware, channel)
	logger.Info("Start gRPC server", zap.String("port", grpcAddress))

	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatal(err.Error())
	}
}
