package wifi

import (
	"errors"
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

var errBadType = errors.New("mock: invalid type assertion")

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	var err error
	if e := args.Error(1); e != nil {
		err = fmt.Errorf("mock error: %w", e)
	}

	if args.Get(0) == nil {
		return nil, err
	}

	result, ok := args.Get(0).([]*wifi.Interface)
	if !ok {
		if err != nil {
			return nil, fmt.Errorf("%w: %w", errBadType, err)
		}
		return nil, errBadType
	}

	return result, err
}

func (m *MockWiFiHandle) AssertExpectations(t mock.TestingT) bool {
	return m.Mock.AssertExpectations(t)
}