package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	mywifi "github.com/task-6/internal/wifi"
)

type MockWiFi struct {
	fn func() ([]*wifi.Interface, error)
}

func (m *MockWiFi) Interfaces() ([]*wifi.Interface, error) { return m.fn() }

func TestWiFiService(t *testing.T) {
	t.Run("GetNames", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			mock := &MockWiFi{fn: func() ([]*wifi.Interface, error) {
				return []*wifi.Interface{{Name: "wlan0"}}, nil
			}}
			svc := mywifi.NewWiFiManager(mock)
			res, err := svc.GetNames()
			assert.NoError(t, err)
			assert.Equal(t, []string{"wlan0"}, res)
		})
		t.Run("Error", func(t *testing.T) {
			mock := &MockWiFi{fn: func() ([]*wifi.Interface, error) {
				return nil, errors.New("fail")
			}}
			svc := mywifi.NewWiFiManager(mock)
			_, err := svc.GetNames()
			assert.Error(t, err)
		})
	})

	t.Run("GetAddresses", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			addr := net.HardwareAddr{0x00, 0x11, 0x22}
			mock := &MockWiFi{fn: func() ([]*wifi.Interface, error) {
				return []*wifi.Interface{{HardwareAddr: addr}}, nil
			}}
			svc := mywifi.NewWiFiManager(mock)
			res, err := svc.GetAddresses()
			assert.NoError(t, err)
			assert.Equal(t, []net.HardwareAddr{addr}, res)
		})
		t.Run("Error", func(t *testing.T) {
			mock := &MockWiFi{fn: func() ([]*wifi.Interface, error) {
				return nil, errors.New("fail")
			}}
			svc := mywifi.NewWiFiManager(mock)
			_, err := svc.GetAddresses()
			assert.Error(t, err)
		})
	})
}