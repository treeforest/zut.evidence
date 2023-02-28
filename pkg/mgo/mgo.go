package mgo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const (
	DefaultAddress     = "localhost:27017"
	DefaultMaxPoolSize = 1024
)

func New(addr string, maxPoolSize uint64) *mongo.Client {
	opts := options.Client()
	opts.SetMaxPoolSize(maxPoolSize)
	opts.ApplyURI(fmt.Sprintf("mongodb://%s", addr))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		panic(err)
	}
	return client
}
