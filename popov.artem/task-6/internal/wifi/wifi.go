package wifi

import (
	"fmt"
	"net"

	"github.com/mdlayher/wifi"
)

type WiFiInterface interface {
	Interfaces() ([]*wifi.Interface, error)
}

type NetworkService struct {
	WiFi WiFiInterface
}

func NewNetworkService(wifi WiFiInterface) NetworkService {
	return NetworkService{WiFi: wifi}
}

func (svc NetworkService) RetrieveMACAddresses() ([]net.HardwareAddr, error) {
	interfaces, err := svc.WiFi.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve interfaces: %w", err)
	}

	macs := make([]net.HardwareAddr, 0, len(interfaces))

	for _, iface := range interfaces {
		macs = append(macs, iface.HardwareAddr)
	}

	return macs, nil
}

func (svc NetworkService) RetrieveInterfaceNames() ([]string, error) {
	interfaces, err := svc.WiFi.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve interfaces: %w", err)
	}

	names := make([]string, 0, len(interfaces))

	for _, iface := range interfaces {
		names = append(names, iface.Name)
	}

	return names, nil
}
