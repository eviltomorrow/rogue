package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//
var (
	DSN     string
	MaxOpen uint64 = 10
	DB      *mongo.Client
)

//
var (
	DefaultConnectTimeout = 10 * time.Second
)

func Build() error {
	pool, err := build(DSN)
	if err != nil {
		return err
	}
	DB = pool
	return err
}

func Close() error {
	if DB == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultConnectTimeout)
	defer cancel()

	return DB.Disconnect(ctx)
}

func build(dsn string) (*mongo.Client, error) {
	if dsn == "" {
		return nil, fmt.Errorf("DSN no set")
	}

	client, err := mongo.NewClient(
		options.Client().ApplyURI(dsn).SetMaxPoolSize(MaxOpen),
	)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultConnectTimeout)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}
