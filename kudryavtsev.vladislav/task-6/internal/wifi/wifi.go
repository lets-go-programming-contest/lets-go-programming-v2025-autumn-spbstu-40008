package wifi

import (
	"fmt"
	"net"

	"github.com/mdlayher/wifi"
)

// WiFiHandle - интерфейс для работы с wifi клиентом.
type WiFiHandle interface {
	Interfaces() ([]*wifi.Interface, error)
}

type WiFiService struct {
	Client WiFiHandle
}

func New(w WiFiHandle) WiFiService {
	return WiFiService{Client: w}
}

// GetAddresses возвращает список MAC-адресов интерфейсов.
func (s WiFiService) GetAddresses() ([]net.HardwareAddr, error) {
	ifaces, err := s.Client.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("interface retrieval failed: %w", err)
	}

	// Аллоцируем слайс сразу нужной длины для производительности
	addrs := make([]net.HardwareAddr, 0, len(ifaces))
	for _, iface := range ifaces {
		addrs = append(addrs, iface.HardwareAddr)
	}

	return addrs, nil
}

// GetNames возвращает список имен интерфейсов (добавлено для полноты покрытия).
func (s WiFiService) GetNames() ([]string, error) {
	ifaces, err := s.Client.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("interface retrieval failed: %w", err)
	}

	names := make([]string, 0, len(ifaces))
	for _, iface := range ifaces {
		names = append(names, iface.Name)
	}

	return names, nil
}