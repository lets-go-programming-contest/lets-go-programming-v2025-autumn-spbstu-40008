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

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	ifaceSlice, ok := args.Get(0).([]*wifi.Interface)
	if !ok {
		// Если приведение типа не удалось, возвращаем ошибку
		return nil, args.Error(1)
	}

	return ifaceSlice, args.Error(1)
}

func (m *MockWiFiHandle) AssertExpectations(t mock.TestingT) bool {
	return m.Mock.AssertExpectations(t)
}
