package wifi_test

import (
	"errors"
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
	t.Run("GetNames_Success", func(t *testing.T) {
		mock := &MockWiFi{fn: func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{{Name: "wlan0"}}, nil
		}}
		svc := mywifi.NewWiFiManager(mock) 
		res, err := svc.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"wlan0"}, res)
	})

	t.Run("GetAddresses_Error", func(t *testing.T) {
		mock := &MockWiFi{fn: func() ([]*wifi.Interface, error) {
			return nil, errors.New("fail")
		}}
		svc := mywifi.NewWiFiManager(mock)
		_, err := svc.GetAddresses()
		assert.Error(t, err)
	})
}