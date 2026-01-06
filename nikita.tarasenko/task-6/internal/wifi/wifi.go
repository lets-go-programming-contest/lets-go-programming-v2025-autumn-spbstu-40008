package wifi

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/mdlayher/wifi"
)

var (
	ErrWiFiAccessDenied = errors.New("access to WiFi interfaces denied")
	ErrWiFiDriverFault  = errors.New("WiFi driver reported fault")
)

type InterfaceEnumerator interface {
	Enumerate(callback func(name string, mac net.HardwareAddr)) error
}

type wifiEnumerator struct {
	source func() ([]*wifi.Interface, error)
}

func (e *wifiEnumerator) Enumerate(callback func(name string, mac net.HardwareAddr)) error {
	ifaces, err := e.source()
	if err != nil {
		var opErr *wifi.OpError
		if errors.As(err, &opErr) {
			switch opErr.Err {
			case wifi.ErrNotSupported:
				return fmt.Errorf("%w: %v", ErrWiFiAccessDenied, err)
			default:
				return fmt.Errorf("%w: %v", ErrWiFiDriverFault, err)
			}
		}
		return fmt.Errorf("unexpected interface fetch error: %w", err)
	}

	for _, iface := range ifaces {
		if iface.Name != "" && iface.HardwareAddr != nil {
			callback(iface.Name, iface.HardwareAddr)
		}
	}

	return nil
}

type DeviceCatalog struct {
	enum InterfaceEnumerator
}

func NewDeviceCatalog(handle func() ([]*wifi.Interface, error)) *DeviceCatalog {
	return &DeviceCatalog{
		enum: &wifiEnumerator{source: handle},
	}
}

func (c *DeviceCatalog) ForEachDevice(fn func(name string, mac net.HardwareAddr)) error {
	return c.enum.Enumerate(fn)
}

func (c *DeviceCatalog) ListDeviceNames() ([]string, error) {
	var names []string
	var mu sync.Mutex

	if err := c.ForEachDevice(func(name string, _ net.HardwareAddr) {
		mu.Lock()
		names = append(names, name)
		mu.Unlock()
	}); err != nil {
		return nil, err
	}

	return names, nil
}

func (c *DeviceCatalog) DeviceMACIterator() <-chan net.HardwareAddr {
	ch := make(chan net.HardwareAddr, 10)
	go func() {
		defer close(ch)
		_ = c.ForEachDevice(func(_ string, mac net.HardwareAddr) {
			ch <- mac
		})
	}()
	return ch
}
