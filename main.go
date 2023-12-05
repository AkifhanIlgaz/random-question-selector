package main

import (
	"context"
	"fmt"
	"log"

	"github.com/AkifhanIlgaz/random-question-selector/cfg"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {

}

func main() {
	config, err := cfg.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not read environment variables", err)
	}

	ctx := context.TODO()

	mongoClient, err := connectToMongo(ctx, config.MongoUri)
	if err != nil {
		log.Fatal("Could not connect to mongo", err)
	}

	redisClient, err := connectToRedis(ctx, config.RedisUri)
	if err != nil {
		log.Fatal("Could not connect to redis", err)
	}
}

func connectToMongo(ctx context.Context, uri string) (*mongo.Client, error) {
	mongoConn := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, mongoConn)
	if err != nil {
		return nil, fmt.Errorf("connect to mongo: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("connect to mongo: %w", err)
	}

	fmt.Println("MongoDB successfully connected...")
	return client, nil
}

func connectToRedis(ctx context.Context, uri string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: uri,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("connect to redis: %w", err)
	}

	err := client.Set(ctx, "test", "Welcome to Golang with Redis and MongoDB", 0).Err()
	if err != nil {
		return nil, fmt.Errorf("connect to redis: %w", err)
	}

	fmt.Println("Redis client connected successfully...")

	return client, nil
}
