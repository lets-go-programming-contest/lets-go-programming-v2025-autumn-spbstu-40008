package wifi

import (
	"errors"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
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
	return args.Get(0).([]*wifi.Interface), args.Error(1)
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m := new(MockWiFiHandle)
		s := New(m)
		m.On("Interfaces").Return([]*wifi.Interface{{Name: "wlan0"}}, nil)

		names, err := s.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"wlan0"}, names)
	})

	t.Run("error", func(t *testing.T) {
		m := new(MockWiFiHandle)
		s := New(m)
		m.On("Interfaces").Return(nil, errors.New("hw fail"))

		_, err := s.GetNames()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "getting interfaces")
	})
}
