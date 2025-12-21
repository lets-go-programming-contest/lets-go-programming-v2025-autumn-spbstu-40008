package wifi_test

import (
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type mockWiFiHandle struct {
	mock.Mock
}

func (mockHandle *mockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := mockHandle.Called()

	var list []*wifi.Interface
	if raw := args.Get(0); raw != nil {
		list = raw.([]*wifi.Interface)
	}

	return list, args.Error(1)
}
