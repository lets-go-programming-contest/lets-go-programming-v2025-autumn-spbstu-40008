package wifi_test

import (
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	var ifaces []*wifi.Interface
	if v := args.Get(0); v != nil {
		if val, ok := v.([]*wifi.Interface); ok {
			ifaces = val
		}
	}

	return ifaces, args.Error(1)
}
