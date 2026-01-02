package wifi

import (
	"fmt"
	"net"

	"github.com/mdlayher/wifi"
)

type WiFiHandler interface {
	Interfaces() ([]*wifi.Interface, error)
}

type WiFiService struct {
	WiFi WiFiHandler
}

func New(wifi WiFiHandler) WiFiService {
	return WiFiService{WiFi: wifi}
}

func (service WiFiService) GetAddresses() ([]net.HardwareAddr, error) {
	interfaces, err := service.WiFi.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("getting interfaces: %w", err)
	}

	addrs := make([]net.HardwareAddr, 0, len(interfaces))

	for _, iface := range interfaces {
		addrs = append(addrs, iface.HardwareAddr)
	}

	return addrs, nil
}

func (service WiFiService) GetNames() ([]string, error) {
	interfaces, err := service.WiFi.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("getting interfaces: %w", err)
	}

	names := make([]string, 0, len(interfaces))

	for _, iface := range interfaces {
		names = append(names, iface.Name)
	}

	return names, nil
}
