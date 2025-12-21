package wifi

import (
    "fmt"
    "net"
    
    "github.com/mdlayher/wifi"
)

type WiFiHandle interface {
    Interfaces() ([]*wifi.Interface, error)
    StationInfo(ifi *wifi.Interface) (*wifi.StationInfo, error)
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
        return nil, fmt.Errorf("get interfaces: %w", err)
    }
    
    var addrs []net.HardwareAddr
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
        return nil, fmt.Errorf("get interfaces: %w", err)
    }
    
    var names []string
    for _, iface := range interfaces {
        names = append(names, iface.Name)
    }
    
    return names, nil
}

func (service WiFiService) GetStationInfo(interfaceName string) (*wifi.StationInfo, error) {
    interfaces, err := service.WiFi.Interfaces()
    if err != nil {
        return nil, fmt.Errorf("get interfaces: %w", err)
    }
    
    for _, iface := range interfaces {
        if iface.Name == interfaceName {
            return service.WiFi.StationInfo(iface)
        }
    }
    
    return nil, fmt.Errorf("interface %s not found", interfaceName)
}
