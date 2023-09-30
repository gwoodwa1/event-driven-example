package main

import (
	"errors"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/stretchr/testify/assert"
)

const topic = "transceiver_telemetry"

func Test(t *testing.T) {
	logger := watermill.NewStdLogger(false, false)

	pubSub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	alert := &Alert{}

	go ConsumeMessages(pubSub, alert)

	time.Sleep(1 * time.Second)

	publishLowPower := func() {
		messageLow := message.NewMessage(watermill.NewUUID(), []byte("-45.0"))
		err := pubSub.Publish(topic, messageLow)
		assert.NoError(t, err)
		time.Sleep(100 * time.Millisecond)
	}

	publishGoodPower := func() {
		messageGood := message.NewMessage(watermill.NewUUID(), []byte("-35.0"))
		err := pubSub.Publish(topic, messageGood)
		assert.NoError(t, err)
		time.Sleep(100 * time.Millisecond)
	}

	publishLowPower()
	assert.True(t, alert.enabled, "alert should be enabled due to low power")

	publishGoodPower()
	assert.False(t, alert.enabled, "alert should be disabled due to good power")

	alert.returnedErr = errors.New("error")

	publishLowPower()
	assert.False(t, alert.enabled, "alert should not be enabled due to error")

	alert.returnedErr = nil
	time.Sleep(100 * time.Millisecond)
	assert.True(t, alert.enabled, "alert should be enabled after error is resolved")

	publishGoodPower()
	assert.False(t, alert.enabled, "alert should be disabled due to good power")
}

type Alert struct {
	enabled     bool
	returnedErr error
}

func (a *Alert) StartAlarm() error {
	if a.returnedErr == nil {
		a.enabled = true
		return nil
	}

	return a.returnedErr
}

func (a *Alert) StopAlarm() error {
	if a.returnedErr == nil {
		a.enabled = false
		return nil
	}

	return a.returnedErr
}
