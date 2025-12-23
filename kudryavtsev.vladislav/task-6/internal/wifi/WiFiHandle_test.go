package wifi

import (
	"errors"
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

var errInvalidMockType = errors.New("mock: return value has invalid type")

// MockWiFiHandle мок-реализация интерфейса WiFiHandle.
type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	// Обработка ошибки (второй аргумент)
	var err error
	if e := args.Error(1); e != nil {
		err = fmt.Errorf("mock error: %w", e)
	}

	// Если первый аргумент nil, возвращаем ошибку сразу
	if args.Get(0) == nil {
		return nil, err
	}

	// Безопасное приведение типа
	ifaces, ok := args.Get(0).([]*wifi.Interface)
	if !ok {
		// Если приведение не удалось, возвращаем спец. ошибку
		if err != nil {
			return nil, fmt.Errorf("%w: original error: %v", errInvalidMockType, err)
		}
		return nil, errInvalidMockType
	}

	return ifaces, err
}

func (m *MockWiFiHandle) AssertExpectations(t mock.TestingT) bool {
	return m.Mock.AssertExpectations(t)
}