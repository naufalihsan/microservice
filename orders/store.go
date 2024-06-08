package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/naufalihsan/msvc-common/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	DbName   = "orders"
	CollName = "orders"
)

func ConnectMongo(user, pass, host, port string) (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", user, pass, host, port)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	return client, err
}

type Store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *Store {
	return &Store{db}
}

func (s *Store) Create(ctx context.Context, order Order) (primitive.ObjectID, error) {
	col := s.db.Database(DbName).Collection(CollName)
	newOrder, err := col.InsertOne(ctx, order)

	return newOrder.InsertedID.(primitive.ObjectID), err
}

func (s *Store) Get(ctx context.Context, customerId, orderId string) (*Order, error) {
	col := s.db.Database(DbName).Collection(CollName)
	objOrderId, _ := primitive.ObjectIDFromHex(orderId)

	var order Order
	err := col.FindOne(ctx, bson.M{"_id": objOrderId, "customerId": customerId}).Decode(&order)

	return &order, err
}

func (s *Store) Update(ctx context.Context, orderId string, order *pb.Order) error {
	col := s.db.Database(DbName).Collection(CollName)
	objOrderId, _ := primitive.ObjectIDFromHex(orderId)

	_, err := col.UpdateOne(ctx,
		bson.M{"_id": objOrderId},
		bson.M{"$set": bson.M{
			"paymentLink": order.PaymentLink,
			"status":      order.Status,
		}})

	return err
}
