package main

import (
	"fmt"
	"math/rand"
	"time"
)

// DataType represents the type of telemetry data.
type DataType int

// Define constants for each data type.
const (
	BroadcastsPkts DataType = iota
	InputDrops
	CrcErrors
)

// Mapping of DataType to its string representation.
var dataTypeStr = map[DataType]string{
	BroadcastsPkts: "broadcasts_pkts",
	InputDrops:     "input_drops",
	CrcErrors:      "crc_errors",
}

// TelemetryTask represents a task to collect a certain type of telemetry data from a device.
type TelemetryTask struct {
	DeviceID string
	DataType DataType
}

// Generates a random value for the given dataType.
func randomValueForDataType(dataType DataType) int {
	switch dataType {
	case BroadcastsPkts:
		return 5000 + rand.Intn(10001)
	case InputDrops:
		return rand.Intn(501)
	case CrcErrors:
		return rand.Intn(101)
	default:
		return 0
	}
}

// TelemetryCollectorClient simulates a client to collect telemetry data.
type TelemetryCollectorClient struct{}

// Collects data from a device. Simulates random failures.
func (c *TelemetryCollectorClient) CollectData(deviceID string, dataType DataType) error {
	// Simulate a 10% chance of failure.
	if rand.Intn(10) == 0 {
		return fmt.Errorf("failed to collect data from device %s", deviceID)
	}

	// Generate and print random telemetry data.
	value := randomValueForDataType(dataType)
	fmt.Printf("Collected %s: %d from device %s\n", dataTypeStr[dataType], value, deviceID)
	return nil
}

// TelemetryWorker is responsible for managing the collection of telemetry data tasks.
type TelemetryWorker struct {
	queue           chan TelemetryTask
	collectorClient *TelemetryCollectorClient
}

// Creates a new telemetry worker.
func NewTelemetryWorker(client *TelemetryCollectorClient) *TelemetryWorker {
	return &TelemetryWorker{
		queue:           make(chan TelemetryTask, 100),
		collectorClient: client,
	}
}

// Sends a task to the worker's queue.
func (w *TelemetryWorker) Send(task TelemetryTask) {
	w.queue <- task
}

// Continuously processes tasks from the queue.
func (w *TelemetryWorker) Run() {
	for task := range w.queue {
		err := w.collectorClient.CollectData(task.DeviceID, task.DataType)
		if err != nil {
			fmt.Printf("Error: %s. Retrying...\n", err)
			w.Send(task)
		}
	}
}

func main() {
	client := &TelemetryCollectorClient{}
	worker := NewTelemetryWorker(client)
	go worker.Run()

	devices := []string{"dc-router-1", "dc-router-2", "dc-router-3"}
	dataTypes := []DataType{BroadcastsPkts, InputDrops, CrcErrors}

	// Continuously send tasks to collect telemetry data every second.
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			for _, device := range devices {
				for _, dataType := range dataTypes {
					worker.Send(TelemetryTask{DeviceID: device, DataType: dataType})
				}
			}
		}
	}
}
