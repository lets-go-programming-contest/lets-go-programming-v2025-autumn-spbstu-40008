package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	wifipkg "task-6/internal/wifi"
)

var errTestInterfaces = errors.New("test interfaces error")

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockWiFi := new(MockWiFiHandler)
		hw, _ := net.ParseMAC("00:11:22:33:44:55")
		ifaces := []*wifi.Interface{{HardwareAddr: hw}}

		mockWiFi.On("Interfaces").Return(ifaces, nil)

		svc := wifipkg.New(mockWiFi)
		res, err := svc.GetAddresses()

		require.NoError(t, err)
		assert.Equal(t, []net.HardwareAddr{hw}, res)
		mockWiFi.AssertExpectations(t)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()

		mockWiFi := new(MockWiFiHandler)
		mockWiFi.On("Interfaces").Return(nil, errTestInterfaces)

		svc := wifipkg.New(mockWiFi)
		res, err := svc.GetAddresses()

		require.Error(t, err)
		assert.Nil(t, res)
		mockWiFi.AssertExpectations(t)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockWiFi := new(MockWiFiHandler)
		ifaces := []*wifi.Interface{{Name: "wlan0"}, {Name: "wlan1"}}

		mockWiFi.On("Interfaces").Return(ifaces, nil)

		svc := wifipkg.New(mockWiFi)
		res, err := svc.GetNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"wlan0", "wlan1"}, res)
		mockWiFi.AssertExpectations(t)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()

		mockWiFi := new(MockWiFiHandler)
		mockWiFi.On("Interfaces").Return(nil, errTestInterfaces)

		svc := wifipkg.New(mockWiFi)
		res, err := svc.GetNames()

		require.Error(t, err)
		assert.Nil(t, res)
		mockWiFi.AssertExpectations(t)
	})
}
