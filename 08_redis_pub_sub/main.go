package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const (
	// set some thresholds which we could use in our logic
	DropsLow      = 500
	DropsMedium   = 1000
	DropsHigh     = 10000
	DropsVeryHigh = 100000
)

type TelemetryData struct {
	Hostname    string `json:"hostname"`
	Interface   string `json:"interface"`
	InputErrors int    `json:"input_errors"`
}

func main() {
	logrus.SetLevel(logrus.InfoLevel)

	// Watermill setup
	logger := watermill.NewStdLogger(false, false)
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		logger.Error("Failed to publish message", err, watermill.LogFields{"topic": "packet-counter-errors", "error": err.Error()})
	}

	// Subscribers
	go processMessages("packet-counter-errors", "packet-error-logs", logPacketErrors)
	go processMessages("packet-counter-errors", "packet-error-alerts", alertOnPacketErrors)

	e := echo.New()

	e.POST("/telemetry-data", func(c echo.Context) error {
		var data TelemetryData
		err := c.Bind(&data)
		if err != nil {
			// Return a 400 Bad Request status code with a relevant error message
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Malformed data received"})
		}

		// Check if the data is valid (you can add more validation checks if needed)
		if data.Hostname == "" || data.InputErrors == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Incomplete or invalid data received"})
		}
		if data.InputErrors >= DropsHigh {
			// Convert the relevant data to a JSON string
			jsonData, _ := json.Marshal(data)
			publisher.Publish("packet-counter-errors", message.NewMessage(watermill.NewUUID(), jsonData))
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "Data received and processed"})
	})

	logrus.Info("Server starting...")

	err = e.Start(":8080")
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func processMessages(topic, consumerGroup string, action func(context.Context, TelemetryData) error) {
	logger := watermill.NewStdLogger(false, false)
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	logger.Info("Starting subscriber for topic", watermill.LogFields{"topic": topic})

	subscriber, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: consumerGroup,
	}, logger)
	if err != nil {
		panic(err)
	}

	messages, err := subscriber.Subscribe(context.Background(), topic)
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		var packetErrorData TelemetryData
		err := json.Unmarshal(msg.Payload, &packetErrorData)
		if err != nil {
			logger.Error("Failed to deserialize message payload", err, watermill.LogFields{"topic": topic})
			msg.Nack()
			continue
		}

		logger.Info("Received message for topic", watermill.LogFields{"topic": topic, "hostname": packetErrorData.Hostname, "interface": packetErrorData.Interface, "input_errors": packetErrorData.InputErrors})

		err = action(msg.Context(), packetErrorData)
		if err != nil {
			msg.Nack()
		} else {
			msg.Ack()
		}
	}
}

func logPacketErrors(ctx context.Context, data TelemetryData) error {
	logrus.Warnf("High Input errors detected on router: %s on interface %s with %d errors", data.Hostname, data.Interface, data.InputErrors)
	return nil
}

func alertOnPacketErrors(ctx context.Context, data TelemetryData) error {
	// Here, you can integrate with any alerting system you have in place.
	logrus.Errorf("ALERT! High Input errors on router: %s on interface %s with %d errors", data.Hostname, data.Interface, data.InputErrors)
	return nil
}
