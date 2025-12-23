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
	Client WiFiHandle
}

func New(w WiFiHandle) WiFiService {
	return WiFiService{Client: w}
}

func (s WiFiService) GetAddresses() ([]net.HardwareAddr, error) {
	ifaces, err := s.Client.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("fetch interfaces: %w", err)
	}

	result := make([]net.HardwareAddr, 0, len(ifaces))
	for _, i := range ifaces {
		result = append(result, i.HardwareAddr)
	}

	return result, nil
}

func (s WiFiService) GetNames() ([]string, error) {
	ifaces, err := s.Client.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("fetch interfaces: %w", err)
	}

	result := make([]string, 0, len(ifaces))
	for _, i := range ifaces {
		result = append(result, i.Name)
	}

	return result, nil
}
