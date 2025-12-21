package wifi

import (
	"errors"
	"fmt"
	"net"

	"github.com/mdlayher/wifi"
)

var (
	ErrInterfaceFetch    = errors.New("failed to fetch interfaces")
	ErrNoValidInterfaces = errors.New("no valid network interfaces found")
)

type InterfaceProvider interface {
	Interfaces() ([]*wifi.Interface, error)
}

type NetworkService struct {
	provider InterfaceProvider
}

func New(provider InterfaceProvider) NetworkService {
	return NetworkService{provider: provider}
}

func (s NetworkService) GetAddresses() ([]net.HardwareAddr, error) {
	interfaces, err := s.provider.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInterfaceFetch, err)
	}

	if len(interfaces) == 0 {
		return nil, fmt.Errorf("%w", ErrNoValidInterfaces)
	}

	addresses := make([]net.HardwareAddr, 0, len(interfaces))

	for _, iface := range interfaces {
		if len(iface.HardwareAddr) > 0 {
			addresses = append(addresses, iface.HardwareAddr)
		}
	}

	if len(addresses) == 0 {
		return nil, fmt.Errorf("%w", ErrNoValidInterfaces)
	}

	return addresses, nil
}

func (s NetworkService) GetNames() ([]string, error) {
	interfaces, err := s.provider.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInterfaceFetch, err)
	}

	if len(interfaces) == 0 {
		return nil, fmt.Errorf("%w", ErrNoValidInterfaces)
	}

	names := make([]string, 0, len(interfaces))

	for _, iface := range interfaces {
		if iface.Name != "" {
			names = append(names, iface.Name)
		}
	}

	if len(names) == 0 {
		return nil, fmt.Errorf("%w", ErrNoValidInterfaces)
	}

	return names, nil
}