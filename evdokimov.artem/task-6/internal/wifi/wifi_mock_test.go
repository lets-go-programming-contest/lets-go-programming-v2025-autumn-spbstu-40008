package wifi_test

import (
	"errors"
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

var errTypeAssertion = errors.New("type assertion failed")

type mockWiFiHandle struct {
	mock.Mock
}

func (m *mockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	var err error
	if args.Error(1) != nil {
		err = fmt.Errorf("mock error: %w", args.Error(1))
	}

	raw := args.Get(0)
	if raw == nil {
		return nil, err
	}

	list, ok := raw.([]*wifi.Interface)
	if !ok {
		if err != nil {
			return nil, fmt.Errorf("%w: %w", errTypeAssertion, err)
		}

		return nil, errTypeAssertion
	}

	return list, err
}
