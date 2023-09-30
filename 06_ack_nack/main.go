package main

import (
	"context"
	"fmt"
	"strconv"
	"github.com/ThreeDotsLabs/watermill/message"
)

type AlarmClient interface {
	StartAlarm() error
	StopAlarm() error
}

func ConsumeMessages(sub message.Subscriber, alarmClient AlarmClient) {
	messages, err := sub.Subscribe(context.Background(), "transceiver_telemetry")
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		powerLevelStr := string(msg.Payload)
		powerLevel, err := strconv.ParseFloat(powerLevelStr, 64)
		if err != nil {
			fmt.Println("Error parsing power level:", err)
			msg.Nack()
			continue
		}

		if powerLevel < -40.0 {
			err = alarmClient.StartAlarm()
		} else {
			err = alarmClient.StopAlarm()
		}

		if err == nil {
			msg.Ack()
		} else {
			msg.Nack()
		}
	}
}
