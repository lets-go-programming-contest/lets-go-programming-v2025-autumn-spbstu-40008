package netif

import (
	"fmt"
	"net"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockInterfaceHandler struct {
	mock.Mock
}

func (m *MockInterfaceHandler) FetchInterfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	interfaces, ok := args.Get(0).([]*wifi.Interface)
	if !ok {
		return nil, fmt.Errorf("type assertion failed for interfaces slice")
	}
	
	return interfaces, args.Error(1)
}

func (m *MockInterfaceHandler) ValidateExpectations(t mock.TestingT) bool {
	return m.Mock.AssertExpectations(t)
}

func createTestInterfaceData(name, macAddress string) *wifi.Interface {
	mac, _ := net.ParseMAC(macAddress)
	return &wifi.Interface{
		Name:         name,
		HardwareAddr: mac,
	}
}