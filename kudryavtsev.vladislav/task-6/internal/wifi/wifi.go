package wifi

import (
	"fmt"
	"net"

	"github.com/mdlayher/wifi"
)

type WiFi interface {
	Interfaces() ([]*wifi.Interface, error)
}

type WiFiService struct {
	Client WiFi
}

func New(w WiFi) WiFiService {
	return WiFiService{Client: w}
}

func (s WiFiService) GetAddresses() ([]net.HardwareAddr, error) {
	ifaces, err := s.Client.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch interfaces: %w", err)
	}

	addrs := make([]net.HardwareAddr, 0, len(ifaces))
	for _, iface := range ifaces {
		addrs = append(addrs, iface.HardwareAddr)
	}

	return addrs, nil
}

func (s WiFiService) GetNames() ([]string, error) {
	ifaces, err := s.Client.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch interfaces: %w", err)
	}

	names := make([]string, 0, len(ifaces))
	for _, iface := range ifaces {
		names = append(names, iface.Name)
	}

	return names, nil
}