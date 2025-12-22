package wifi

import (
	"fmt"

	"github.com/mdlayher/wifi"
)

type WiFiHandle interface {
	Interfaces() ([]*wifi.Interface, error)
}

type NetManager struct {
	WiFi WiFiHandle
}

func New(h WiFiHandle) NetManager {
	return NetManager{WiFi: h}
}

func (m NetManager) GetActiveInterfaces() ([]string, error) {
	ifaces, err := m.WiFi.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("fetch wifi interfaces: %w", err)
	}

	names := make([]string, 0, len(ifaces))
	for _, iface := range ifaces {
		names = append(names, iface.Name)
	}
	return names, nil
}
