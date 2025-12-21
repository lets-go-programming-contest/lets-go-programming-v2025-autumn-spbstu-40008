package wifi_test

import (
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

// Ошибка для type assertion
var ErrTypeAssertionFailed = fmt.Errorf("type assertion failed")

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	var err error
	if args.Error(1) != nil {
		err = fmt.Errorf("mock error: %w", args.Error(1))
	}

	if args.Get(0) == nil {
		return nil, err
	}

	ifaceSlice, ok := args.Get(0).([]*wifi.Interface)
	if !ok {
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrTypeAssertionFailed, err)
		}
		return nil, ErrTypeAssertionFailed
	}

	return ifaceSlice, err
}

func (m *MockWiFiHandle) AssertExpectations(t mock.TestingT) bool {
	return m.Mock.AssertExpectations(t)
}
