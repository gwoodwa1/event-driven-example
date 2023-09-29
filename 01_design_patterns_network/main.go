package main

import (
	"fmt"
	"log"
)

// Device struct represents a network device in the system.
type Device struct {
	IPAddress string
}

// DeviceRepository interface represents a component responsible for fetching device information.
type DeviceRepository interface {
	GetDevice(ipAddress string) (Device, error)
}

// ConfigurationClient interface represents a component responsible for configuring devices.
type ConfigurationClient interface {
	ConfigureDevice(d Device) error
}

// MonitoringClient interface represents a component responsible for monitoring devices.
type MonitoringClient interface {
	MonitorDevice(d Device) error
}

// NetworkHandler struct is used to handle network operations by interacting with the corresponding components.
type NetworkHandler struct {
	deviceRepository   DeviceRepository
	configurationClient ConfigurationClient
	monitoringClient    MonitoringClient
}

// NewNetworkHandler is a constructor for the NetworkHandler struct.
func NewNetworkHandler(
	deviceRepository DeviceRepository,
	configurationClient ConfigurationClient,
	monitoringClient MonitoringClient,
) NetworkHandler {
	return NetworkHandler{
		deviceRepository:   deviceRepository,
		configurationClient: configurationClient,
		monitoringClient:    monitoringClient,
	}
}

// PerformNetworkOperation method is responsible for performing network operations on a device.
func (n NetworkHandler) PerformNetworkOperation(ipAddress string) error {
	device, err := n.deviceRepository.GetDevice(ipAddress)
	if err != nil {
		return fmt.Errorf("failed to get device: %v", err)
	}

	if err := n.configurationClient.ConfigureDevice(device); err != nil {
		return fmt.Errorf("failed to configure device: %v", err)
	}

	if err := n.monitoringClient.MonitorDevice(device); err != nil {
		return fmt.Errorf("failed to monitor device: %v", err)
	}

	return nil
}

// Mock implementations for the interfaces
type MockDeviceRepository struct{}

func (m MockDeviceRepository) GetDevice(ipAddress string) (Device, error) {
	// Mock implementation
	return Device{IPAddress: ipAddress}, nil
}

type MockConfigurationClient struct{}

func (m MockConfigurationClient) ConfigureDevice(d Device) error {
	// Mock implementation
	log.Printf("Configuring device with IP: %s", d.IPAddress)
	return nil
}

type MockMonitoringClient struct{}

func (m MockMonitoringClient) MonitorDevice(d Device) error {
	// Mock implementation
	log.Printf("Monitoring device with IP: %s", d.IPAddress)
	return nil
}

func main() {
	networkHandler := NewNetworkHandler(MockDeviceRepository{}, MockConfigurationClient{}, MockMonitoringClient{})
	ipAddress := "192.168.1.1"

	if err := networkHandler.PerformNetworkOperation(ipAddress); err != nil {
		log.Fatalf("Failed to perform network operation: %v", err)
	}
}
