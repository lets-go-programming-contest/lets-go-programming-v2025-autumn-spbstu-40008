package wifi

import (
	"fmt"
	"net"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockInterfaceSource struct {
	mock.Mock
}

func (m *MockInterfaceSource) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	interfaces, ok := args.Get(0).([]*wifi.Interface)
	if !ok {
		return nil, fmt.Errorf("type assertion failed for interface slice")
	}
	return interfaces, args.Error(1)
}

func createTestInterface(name, macStr string) *wifi.Interface {
	mac, _ := net.ParseMAC(macStr)
	return &wifi.Interface{
		Name: name,
		HardwareAddr: mac,
	}
}
