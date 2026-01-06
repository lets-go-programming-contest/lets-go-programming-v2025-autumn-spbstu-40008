package wifi

import (
	"fmt"
	"net"

	"github.com/mdlayher/wifi"
)

type WiFiConnector interface {
	NetworkInterfaces() ([]*wifi.Interface, error)
}

type WiFiManager struct {
	adapter WiFiConnector
}

func Create(adapter WiFiConnector) WiFiManager {
	return WiFiManager{adapter: adapter}
}

func (manager WiFiManager) FetchHardwareAddresses() ([]net.HardwareAddr, error) {
	interfaces, err := manager.adapter.NetworkInterfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve network interfaces: %w", err)
	}

	addresses := make([]net.HardwareAddr, 0, len(interfaces))

	for _, iface := range interfaces {
		addresses = append(addresses, iface.HardwareAddr)
	}

	return addresses, nil
}

func (manager WiFiManager) FetchInterfaceNames() ([]string, error) {
	interfaces, err := manager.adapter.NetworkInterfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve network interfaces: %w", err)
	}

	names := make([]string, 0, len(interfaces))

	for _, iface := range interfaces {
		names = append(names, iface.Name)
	}

	return names, nil
}
