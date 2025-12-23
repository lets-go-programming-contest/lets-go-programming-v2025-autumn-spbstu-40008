package wifi_test

import (
	"errors"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mywifi "github.com/task-6/internal/wifi" // Исправлено имя модуля
)

// Mock WiFiHandle прямо здесь, чтобы не было проблем с файлами
type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*wifi.Interface), args.Error(1)
}

func TestWiFi(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockWiFi := new(MockWiFiHandle)
		service := mywifi.New(mockWiFi)
		mockWiFi.On("Interfaces").Return([]*wifi.Interface{{Name: "wlan0"}}, nil)

		names, err := service.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"wlan0"}, names)
	})

	t.Run("error", func(t *testing.T) {
		mockWiFi := new(MockWiFiHandle)
		service := mywifi.New(mockWiFi)
		mockWiFi.On("Interfaces").Return([]*wifi.Interface(nil), errors.New("err"))

		names, err := service.GetNames()
		assert.Error(t, err)
		assert.Nil(t, names)
	})
}
