package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	wifipkg "popov.artem/task-6/internal/wifi"
)

var (
	errMACRead = errors.New("mac address read failed")
	errIFRead  = errors.New("interface list read failed")
)

func TestNetworkService_RetrieveMACAddresses(t *testing.T) {
	t.Parallel()

	t.Run("successful_retrieval", func(t *testing.T) {
		t.Parallel()

		m := new(wifipkg.MockWiFiHandle)
		hw, _ := net.ParseMAC("00:11:22:33:44:55")
		ifaces := []*wifi.Interface{{HardwareAddr: hw}}

		m.On("Interfaces").Return(ifaces, nil)

		svc := wifipkg.NewNetworkService(m)
		res, err := svc.RetrieveMACAddresses()

		require.NoError(t, err)
		assert.Equal(t, []net.HardwareAddr{hw}, res)
		m.AssertExpectations(t)
	})

	t.Run("interface_error", func(t *testing.T) {
		t.Parallel()

		m := new(wifipkg.MockWiFiHandle)
		m.On("Interfaces").Return(nil, errMACRead)

		svc := wifipkg.NewNetworkService(m)
		res, err := svc.RetrieveMACAddresses()

		require.Error(t, err)
		assert.Nil(t, res)
		m.AssertExpectations(t)
	})
}

func TestNetworkService_RetrieveInterfaceNames(t *testing.T) {
	t.Parallel()

	t.Run("successful_names", func(t *testing.T) {
		t.Parallel()

		m := new(wifipkg.MockWiFiHandle)
		ifaces := []*wifi.Interface{{Name: "wlan0"}, {Name: "wlan1"}}

		m.On("Interfaces").Return(ifaces, nil)

		svc := wifipkg.NewNetworkService(m)
		res, err := svc.RetrieveInterfaceNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"wlan0", "wlan1"}, res)
		m.AssertExpectations(t)
	})

	t.Run("interface_list_error", func(t *testing.T) {
		t.Parallel()

		m := new(wifipkg.MockWiFiHandle)
		m.On("Interfaces").Return(nil, errIFRead)

		svc := wifipkg.NewNetworkService(m)
		res, err := svc.RetrieveInterfaceNames()

		require.Error(t, err)
		assert.Nil(t, res)
		m.AssertExpectations(t)
	})
}
