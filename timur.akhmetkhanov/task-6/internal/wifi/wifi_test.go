package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	internalWifi "task-6/internal/wifi"
)

var (
	errDriver = errors.New("driver error")
	errFail   = errors.New("fail")
)

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	hwAddr1, _ := net.ParseMAC("00:00:5e:00:53:01")
	hwAddr2, _ := net.ParseMAC("00:00:5e:00:53:02")

	mockInterfaces := []*wifi.Interface{
		{HardwareAddr: hwAddr1, Name: "wlan0"},
		{HardwareAddr: hwAddr2, Name: "wlan1"},
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockWiFi := new(MockWiFiHandle)
		service := internalWifi.New(mockWiFi)

		mockWiFi.On("Interfaces").Return(mockInterfaces, nil).Once()

		addrs, err := service.GetAddresses()

		require.NoError(t, err)
		assert.Len(t, addrs, 2)
		assert.Equal(t, hwAddr1, addrs[0])
		assert.Equal(t, hwAddr2, addrs[1])
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		mockWiFi := new(MockWiFiHandle)
		service := internalWifi.New(mockWiFi)

		mockWiFi.On("Interfaces").Return(nil, errDriver).Once()

		addrs, err := service.GetAddresses()

		require.Error(t, err)
		assert.Nil(t, addrs)
		assert.Contains(t, err.Error(), "getting interfaces")
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	mockInterfaces := []*wifi.Interface{
		{Name: "eth0"},
		{Name: "wlan0"},
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockWiFi := new(MockWiFiHandle)
		service := internalWifi.New(mockWiFi)

		mockWiFi.On("Interfaces").Return(mockInterfaces, nil).Once()

		names, err := service.GetNames()

		require.NoError(t, err)
		assert.Len(t, names, 2)
		assert.Equal(t, "eth0", names[0])
		assert.Equal(t, "wlan0", names[1])
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		mockWiFi := new(MockWiFiHandle)
		service := internalWifi.New(mockWiFi)

		mockWiFi.On("Interfaces").Return(nil, errFail).Once()

		names, err := service.GetNames()

		require.Error(t, err)
		assert.Nil(t, names)
	})
}
