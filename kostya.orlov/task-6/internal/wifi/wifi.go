package wifi

import (
	"fmt"
	"net"

	"github.com/mdlayher/wifi"
)

type WiFiHandle interface {
	Interfaces() ([]*wifi.Interface, error)
}

type WiFiService struct {
	WiFi WiFiHandle
}


func NewWiFiManager(w WiFiHandle) WiFiService {
	return WiFiService{WiFi: w}
}

func (s WiFiService) GetNames() ([]string, error) {
	ifaces, err := s.WiFi.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("getting interfaces: %w", err)
	}

	names := make([]string, 0, len(ifaces))
	for _, iface := range ifaces {
		names = append(names, iface.Name)
	}
	return names, nil
}

func (s WiFiService) GetAddresses() ([]net.HardwareAddr, error) {
	ifaces, err := s.WiFi.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("getting interfaces: %w", err)
	}

	addrs := make([]net.HardwareAddr, 0, len(ifaces))
	for _, iface := range ifaces {
		addrs = append(addrs, iface.HardwareAddr)
	}
	return addrs, nil
}
