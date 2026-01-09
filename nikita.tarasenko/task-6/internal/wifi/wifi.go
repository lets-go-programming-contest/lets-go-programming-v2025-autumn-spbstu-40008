package wifi

import (
	"fmt"
	"net"

	"github.com/mdlayher/wifi"
)

type WiFiHandle interface {
	Interfaces() ([]*wifi.Interface, error)
}

type WiFiService struct {
	wh WiFiHandle
}

func New(wh WiFiHandle) WiFiService {
	return WiFiService{wh: wh}
}

func (ws WiFiService) GetAddresses() ([]net.HardwareAddr, error) {
	ifaces, err := ws.wh.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("не удалось получить интерфейсы: %w", err)
	}

	addresses := make([]net.HardwareAddr, 0, len(ifaces))

	for _, i := range ifaces {
		addresses = append(addresses, i.HardwareAddr)
	}

	return addresses, nil
}

func (ws WiFiService) GetNames() ([]string, error) {
	ifaces, err := ws.wh.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения интерфейсов: %w", err)
	}

	names := make([]string, 0, len(ifaces))

	for _, i := range ifaces {
		names = append(names, i.Name)
	}

	return names, nil
}
