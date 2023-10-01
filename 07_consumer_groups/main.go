package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

type AlertsClient interface {
	SendAlert(telemetryID string) error
}

type LoggingClient interface {
	LogTelemetryData(telemetryID string) error
}

func Subscribe(
	alertsClient AlertsClient,
	loggingClient LoggingClient,
) error {
	logger := watermill.NewStdLogger(false, false)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	// Subscriber for alerts on Packet Counter Errors
	packetErrorAlertSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "packet-error-alerts",
	}, logger)
	if err != nil {
		return err
	}

	// Subscriber for logging Packet Counter Errors
	packetErrorLogSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "packet-error-logs",
	}, logger)
	if err != nil {
		return err
	}

	go processMessages(packetErrorAlertSub, alertsClient.SendAlert, "packet-counter-errors")
	go processMessages(packetErrorLogSub, loggingClient.LogTelemetryData, "packet-counter-errors")

	return nil
}

func processMessages(sub message.Subscriber, action func(telemetryID string) error, topic string) {
	messages, err := sub.Subscribe(context.Background(), topic)
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		telemetryID := string(msg.Payload)

		err := action(telemetryID)
		if err != nil {
			msg.Nack()
		} else {
			msg.Ack()
		}
	}
}

func main() {
	// Mock implementations for AlertsClient and LoggingClient
	alertsClient := &mockAlertsClient{}
	loggingClient := &mockLoggingClient{}

	// Call the Subscribe function
	err := Subscribe(alertsClient, loggingClient)
	if err != nil {
		panic(err)
	}

	// Keep the application running (for demonstration purposes)
	select {}
}

// Mock implementation for AlertsClient
type mockAlertsClient struct{}

func (m *mockAlertsClient) SendAlert(telemetryID string) error {
	// Mock sending an alert
	fmt.Printf("Alert sent for telemetry ID: %s\n", telemetryID)
	return nil
}

// Mock implementation for LoggingClient
type mockLoggingClient struct{}

func (m *mockLoggingClient) LogTelemetryData(telemetryID string) error {
	// Mock logging telemetry data
	fmt.Printf("Telemetry data logged for ID: %s\n", telemetryID)
	return nil
}


