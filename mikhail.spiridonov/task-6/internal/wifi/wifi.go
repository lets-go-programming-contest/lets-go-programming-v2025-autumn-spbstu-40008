package wifi

import (
	"fmt"
	"net"

	mdlayherWifi "github.com/mdlayher/wifi"
)

var (
	ErrInterfaceNotFound = errors.New("interface not found")
	ErrGetInterfaces     = errors.New("failed to get interfaces")
	ErrGetStationInfo    = errors.New("get station info failed")
)

type WiFiHandle interface {
	Interfaces() ([]*mdlayherWifi.Interface, error)
	StationInfo(ifi *mdlayherWifi.Interface) (*mdlayherWifi.StationInfo, error)
}

type WiFiService struct {
	WiFi WiFiHandle
}

func New(wifiHandle WiFiHandle) WiFiService {
	return WiFiService{WiFi: wifiHandle}
}

func (service WiFiService) GetAddresses() ([]net.HardwareAddr, error) {
	interfaces, err := service.WiFi.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetInterfaces, err)
	}

	addrs := make([]net.HardwareAddr, 0, len(interfaces))

	for _, iface := range interfaces {
		if iface.HardwareAddr != nil {
			addrs = append(addrs, iface.HardwareAddr)
		}
	}

	return addrs, nil
}

func (service WiFiService) GetInterfaceNames() ([]string, error) {
	interfaces, err := service.WiFi.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetInterfaces, err)
	}

	names := make([]string, 0, len(interfaces))

	for _, iface := range interfaces {
		names = append(names, iface.Name)
	}

	return names, nil
}

func (service WiFiService) GetStationInfo(interfaceName string) (*mdlayherWifi.StationInfo, error) {
	interfaces, err := service.WiFi.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetInterfaces, err)
	}

	for _, iface := range interfaces {
		if iface.Name == interfaceName {
			info, err := service.WiFi.StationInfo(iface)
			if err != nil {
				return nil, fmt.Errorf("%w: %w", ErrGetStationInfo, err)
			}

			return info, nil
		}
	}

	return nil, fmt.Errorf("interface %s: %w", interfaceName, ErrInterfaceNotFound)
}
