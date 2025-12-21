package netif

import (
	"errors"
	"fmt"
	"net"

	"github.com/mdlayher/wifi"
)

var (
	ErrInterfaceFetch    = errors.New("failed to fetch network interfaces")
	ErrNoValidInterfaces = errors.New("no valid network interfaces found")
	ErrInvalidMACAddress = errors.New("invalid MAC address format detected")
)

type InterfaceHandler interface {
	FetchInterfaces() ([]*wifi.Interface, error)
}

type NetworkService struct {
	handler InterfaceHandler
}

func NewNetworkService(handler InterfaceHandler) *NetworkService {
	return &NetworkService{handler: handler}
}

func (s *NetworkService) GetHardwareAddresses() ([]net.HardwareAddr, error) {
	interfaces, err := s.handler.FetchInterfaces()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInterfaceFetch, err)
	}

	if len(interfaces) == 0 {
		return nil, fmt.Errorf("%w: empty interface list", ErrNoValidInterfaces)
	}

	var macAddresses []net.HardwareAddr
	for _, iface := range interfaces {
		if len(iface.HardwareAddr) == 6 {
			macAddresses = append(macAddresses, iface.HardwareAddr)
		}
	}

	if len(macAddresses) == 0 {
		return nil, fmt.Errorf("%w: no valid MAC addresses detected", ErrInvalidMACAddress)
	}

	return macAddresses, nil
}

func (s *NetworkService) GetInterfaceIdentifiers() ([]string, error) {
	interfaces, err := s.handler.FetchInterfaces()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInterfaceFetch, err)
	}

	if len(interfaces) == 0 {
		return nil, fmt.Errorf("%w: no interfaces available", ErrNoValidInterfaces)
	}

	interfaceNames := make([]string, 0, len(interfaces))
	for _, iface := range interfaces {
		if iface.Name != "" {
			interfaceNames = append(interfaceNames, iface.Name)
		}
	}

	if len(interfaceNames) == 0 {
		return nil, fmt.Errorf("%w: all interface names are empty", ErrNoValidInterfaces)
	}

	return interfaceNames, nil
}