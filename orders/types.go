package main

import (
	"context"

	pb "github.com/naufalihsan/msvc-common/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error)
	ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) ([]*pb.Product, error)
	GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error)
	UpdateOrder(ctx context.Context, order *pb.Order) (*pb.Order, error)
}

type OrderStore interface {
	Create(context.Context, Order) (primitive.ObjectID, error)
	Get(ctx context.Context, id, customerId string) (*Order, error)
	Update(ctx context.Context, orderId string, order *pb.Order) error
}

type Order struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	CustomerId  string             `bson:"customerId,omitempty"`
	Status      string             `bson:"status,omitempty"`
	PaymentLink string             `bson:"paymentLink,omitempty"`
	Products    []*pb.Product      `bson:"products,omitempty"`
}

func (o *Order) ToProto() *pb.Order {
	return &pb.Order{
		Id:          o.Id.Hex(),
		CustomerId:  o.CustomerId,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
	}
}
