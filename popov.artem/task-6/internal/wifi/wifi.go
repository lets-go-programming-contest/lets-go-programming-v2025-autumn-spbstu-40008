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
	WiFi WiFi
}

func New(w WiFi) WiFiService {
	return WiFiService{WiFi: w}
}

func (s WiFiService) GetAddresses() ([]net.HardwareAddr, error) {
	devList, err := s.WiFi.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("interface enumeration failed: %w", err)
	}

	hwList := make([]net.HardwareAddr, 0, len(devList))
	for _, dev := range devList {
		hwList = append(hwList, dev.HardwareAddr)
	}

	return hwList, nil
}

func (s WiFiService) GetNames() ([]string, error) {
	devList, err := s.WiFi.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to list interface names: %w", err)
	}

	nameList := make([]string, 0, len(devList))
	for _, dev := range devList {
		nameList = append(nameList, dev.Name)
	}

	return nameList, nil
}
