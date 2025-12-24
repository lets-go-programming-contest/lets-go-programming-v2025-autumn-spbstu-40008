package wifi_test

import (
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	var ifaces []*wifi.Interface

	v := args.Get(0)
	if v != nil {
		if val, ok := v.([]*wifi.Interface); ok {
			ifaces = val
		}
	}

	if err := args.Error(1); err != nil {
		return ifaces, fmt.Errorf("mock error: %w", err)
	}

	return ifaces, nil
}
