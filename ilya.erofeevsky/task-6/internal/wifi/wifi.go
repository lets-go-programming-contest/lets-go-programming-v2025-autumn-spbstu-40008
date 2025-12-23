package wifi

import (
	"fmt"

	"github.com/mdlayher/wifi"
)

type WiFiHandle interface {
	Interfaces() ([]*wifi.Interface, error)
}

type WiFiService struct {
	WiFi WiFiHandle
}

func New(w WiFiHandle) WiFiService {
	return WiFiService{WiFi: w}
}

func (s WiFiService) GetNames() ([]string, error) {
	ifaces, err := s.WiFi.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("getting interfaces: %w", err)
	}
	names := make([]string, 0, len(ifaces))
	for _, i := range ifaces {
		names = append(names, i.Name)
	}
	return names, nil
}
