package wifi

import (
	"errors"
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

var errBadTypeAssertion = errors.New("mock: type assertion failed")

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	var err error
	if e := args.Error(1); e != nil {
		err = fmt.Errorf("mock error: %w", e)
	}

	result := args.Get(0)
	if result == nil {
		return nil, err
	}

	ifaces, ok := result.([]*wifi.Interface)
	if !ok {
		if err != nil {
			return nil, fmt.Errorf("%w: %w", errBadTypeAssertion, err)
		}
		return nil, errBadTypeAssertion
	}

	return ifaces, err
}

func (m *MockWiFiHandle) AssertExpectations(t mock.TestingT) bool {
	return m.Mock.AssertExpectations(t)
}