package wifi_test

import (
	"errors"
	"net"
	"testing"

	myWifi "github.com/deonik3/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

//go:generate mockery --all --testonly --quiet --outpkg wifi_test --output .

var ErrExpected = errors.New("expected error")

func TestGetAddresses_Success(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := myWifi.New(mockWiFi)

	ifaces := []*wifi.Interface{
		{Name: "wlan0", HardwareAddr: parseMAC(t, "00:11:22:33:44:55")},
		{Name: "wlan1", HardwareAddr: parseMAC(t, "aa:bb:cc:dd:ee:ff")},
	}
	mockWiFi.On("Interfaces").Return(ifaces, nil)

	addrs, err := service.GetAddresses()
	require.NoError(t, err)
	require.Equal(t, []net.HardwareAddr{
		parseMAC(t, "00:11:22:33:44:55"),
		parseMAC(t, "aa:bb:cc:dd:ee:ff"),
	}, addrs)
}

func TestGetAddresses_Error(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := myWifi.New(mockWiFi)

	mockWiFi.On("Interfaces").Return(nil, ErrExpected)

	addrs, err := service.GetAddresses()
	require.ErrorContains(t, err, "getting interfaces")
	require.Nil(t, addrs)
}

func TestGetNames_Success(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := myWifi.New(mockWiFi)

	ifaces := []*wifi.Interface{
		{Name: "wlan0", HardwareAddr: parseMAC(t, "00:11:22:33:44:55")},
		{Name: "eth1", HardwareAddr: parseMAC(t, "11:22:33:44:55:66")},
	}
	mockWiFi.On("Interfaces").Return(ifaces, nil)

	names, err := service.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"wlan0", "eth1"}, names)
}

func TestGetNames_Error(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)
	service := myWifi.New(mockWiFi)

	mockWiFi.On("Interfaces").Return(nil, ErrExpected)

	names, err := service.GetNames()
	require.ErrorContains(t, err, "getting interfaces")
	require.Nil(t, names)
}

func parseMAC(t *testing.T, s string) net.HardwareAddr {
	t.Helper()

	hwAddr, err := net.ParseMAC(s)
	require.NoError(t, err)

	return hwAddr
}
