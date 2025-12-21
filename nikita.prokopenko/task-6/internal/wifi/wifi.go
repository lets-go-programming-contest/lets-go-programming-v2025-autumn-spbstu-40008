package wifi

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/mdlayher/wifi"
)

var (
	ErrInterfaceFetch = errors.New("failed to fetch interfaces")
	ErrNoValidData    = errors.New("no valid interface data")
)

const macAddressLen = 6

type InterfaceSource interface {
	Interfaces() ([]*wifi.Interface, error)
}

type NetworkManager struct {
	source InterfaceSource
}

func CreateManager(source InterfaceSource) *NetworkManager {
	return &NetworkManager{source: source}
}

func (m *NetworkManager) GetMACAddresses() ([]net.HardwareAddr, error) {
	interfaces, err := m.source.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInterfaceFetch, err)
	}
	if len(interfaces) == 0 {
		return nil, fmt.Errorf("%w: empty interface list", ErrNoValidData)
	}
	var macs []net.HardwareAddr
	for _, iface := range interfaces {
		if len(iface.HardwareAddr) == macAddressLen {
			macs = append(macs, iface.HardwareAddr)
		}
	}
	if len(macs) == 0 {
		return nil, fmt.Errorf("%w: no valid MAC addresses", ErrNoValidData)
	}
	return macs, nil
}

func (m *NetworkManager) GetInterfaceNames() ([]string, error) {
	interfaces, err := m.source.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInterfaceFetch, err)
	}
	if len(interfaces) == 0 {
		return nil, fmt.Errorf("%w: no interfaces available", ErrNoValidData)
	}
	names := make([]string, 0, len(interfaces))
	for _, iface := range interfaces {
		if strings.TrimSpace(iface.Name) != "" {
			names = append(names, iface.Name)
		}
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("%w: all names empty", ErrNoValidData)
	}
	return names, nil
}
