package wifi_test

import (
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockWiFiHandler struct {
	mock.Mock
}

func (m *MockWiFiHandler) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	var ifaces []*wifi.Interface
	val := args.Get(0)

	if val != nil {
		var ok bool
		ifaces, ok = val.([]*wifi.Interface)

		if !ok {
			return nil, fmt.Errorf("type assertion failed: %w", args.Error(1))
		}
	}

	if err := args.Error(1); err != nil {
		return ifaces, fmt.Errorf("mock error: %w", err)
	}

	return ifaces, nil
}
