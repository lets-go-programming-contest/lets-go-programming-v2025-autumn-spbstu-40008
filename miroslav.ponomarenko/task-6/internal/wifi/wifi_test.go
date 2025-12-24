package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	internalwifi "rabbitdfs/task-6/internal/wifi"
)

var (
	errHardware = errors.New("hardware error")
	errDriver   = errors.New("driver error")
)

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := new(MockWiFiHandle)
		hw, _ := net.ParseMAC("00:11:22:33:44:55")
		ifaces := []*wifi.Interface{{HardwareAddr: hw}}

		m.On("Interfaces").Return(ifaces, nil)

		svc := internalwifi.New(m)
		res, err := svc.GetAddresses()

		require.NoError(t, err)
		assert.Equal(t, []net.HardwareAddr{hw}, res)
		m.AssertExpectations(t)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()

		m := new(MockWiFiHandle)
		m.On("Interfaces").Return(nil, errHardware)

		svc := internalwifi.New(m)
		res, err := svc.GetAddresses()

		require.Error(t, err)
		assert.Nil(t, res)
		m.AssertExpectations(t)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := new(MockWiFiHandle)
		ifaces := []*wifi.Interface{{Name: "wlan0"}, {Name: "wlan1"}}

		m.On("Interfaces").Return(ifaces, nil)

		svc := internalwifi.New(m)
		res, err := svc.GetNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"wlan0", "wlan1"}, res)
		m.AssertExpectations(t)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()

		m := new(MockWiFiHandle)
		m.On("Interfaces").Return(nil, errDriver)

		svc := internalwifi.New(m)
		res, err := svc.GetNames()

		require.Error(t, err)
		assert.Nil(t, res)
		m.AssertExpectations(t)
	})
}
