package wifi

import (
	"fmt"
	"net"

	"github.com/mdlayher/wifi"
)

type InterfaceProvider interface {
	FetchInterfaces() ([]*wifi.Interface, error)
}

type NetworkService struct {
	Provider InterfaceProvider
}

func NewNetworkService(provider InterfaceProvider) NetworkService {
	return NetworkService{Provider: provider}
}

func (ns NetworkService) RetrieveMACAddresses() ([]net.HardwareAddr, error) {
	devices, err := ns.Provider.FetchInterfaces()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve interfaces: %w", err)
	}

	addresses := make([]net.HardwareAddr, 0, len(devices))
	for _, dev := range devices {
		addresses = append(addresses, dev.HardwareAddr)
	}
	return addresses, nil
}

func (ns NetworkService) RetrieveInterfaceNames() ([]string, error) {
	devices, err := ns.Provider.FetchInterfaces()
	if err != nil {
		return nil, fmt.Errorf("error fetching interface list: %w", err)
	}

	names := make([]string, 0, len(devices))
	for _, dev := range devices {
		names = append(names, dev.Name)
	}
	return names, nil
}
