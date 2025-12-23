package wifi_test

import (
	"errors"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mywifi "github.com/Ilya-Er0fick/task-6/internal/wifi"
)

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

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		m := new(MockWiFiHandle)
		s := mywifi.New(m)
		m.On("Interfaces").Return([]*wifi.Interface{{Name: "wlan0"}}, nil)

		names, err := s.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"wlan0"}, names)
	})

	t.Run("success multiple interfaces", func(t *testing.T) {
		m := new(MockWiFiHandle)
		s := mywifi.New(m)
		m.On("Interfaces").Return([]*wifi.Interface{
			{Name: "wlan0"},
			{Name: "wlan1"},
			{Name: "eth0"},
		}, nil)

		names, err := s.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"wlan0", "wlan1", "eth0"}, names)
	})

	t.Run("success empty interfaces", func(t *testing.T) {
		m := new(MockWiFiHandle)
		s := mywifi.New(m)
		m.On("Interfaces").Return([]*wifi.Interface{}, nil)

		names, err := s.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{}, names)
	})

	t.Run("error", func(t *testing.T) {
		m := new(MockWiFiHandle)
		s := mywifi.New(m)
		m.On("Interfaces").Return(nil, errors.New("hw fail"))

		_, err := s.GetNames()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "getting interfaces")
	})
}
