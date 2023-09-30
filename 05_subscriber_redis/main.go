package main

import (
	"context"
	"fmt"
	"os"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Setup logger
	logger := watermill.NewStdLogger(false, false)

	// Setup Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	// Create a subscriber for Redis Streams
	subscriber, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}
	defer subscriber.Close()

	// Subscribe to the "progress" topic
	messages, err := subscriber.Subscribe(context.Background(), "progress")
	if err != nil {
		panic(err)
	}

	// Process incoming messages
	for msg := range messages {
		fmt.Printf("Message ID: %s - %s%%\n", msg.UUID, msg.Payload)
		msg.Ack()
	}
}
