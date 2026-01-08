package netinfo

import (
    "fmt"
    "github.com/mdlayher/wifi"
)

type Scanner interface {
    Interfaces() ([]*wifi.Interface, error)
}

type WiFiManager struct {
    scanner Scanner
}

func NewWiFiManager(s Scanner) *WiFiManager {
    return &WiFiManager{scanner: s}
}

func (m *WiFiManager) FetchMACAddresses() ([]string, error) {
    ifaces, err := m.scanner.Interfaces()
    if err != nil {
        return nil, fmt.Errorf("network error: %w", err)
    }

    macs := make([]string, 0, len(ifaces))
    for _, i := range ifaces {
        macs = append(macs, i.HardwareAddr.String())
    }
    return macs, nil
}