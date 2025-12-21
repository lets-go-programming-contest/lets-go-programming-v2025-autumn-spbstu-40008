package wifi_test

import (
	"errors"
	"fmt"
	"net"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

var errTypeAssertion = errors.New("type assertion failed for interface slice")

type MockInterfaceSource struct {
	mock.Mock
}

func (m *MockInterfaceSource) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	if args.Get(0) == nil {
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("mock error: %w", err)
		}

		return nil, nil
	}

	interfaces, ok := args.Get(0).([]*wifi.Interface)
	if !ok {
		return nil, errTypeAssertion
	}

	if err := args.Error(1); err != nil {
		return interfaces, fmt.Errorf("mock error: %w", err)
	}

	return interfaces, nil
}

func createTestInterface(name, macStr string) *wifi.Interface {
	mac, _ := net.ParseMAC(macStr)

	return &wifi.Interface{
		Name:         name,
		HardwareAddr: mac,
	}
}
