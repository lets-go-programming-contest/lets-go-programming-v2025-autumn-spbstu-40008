package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	mywifi "github.com/kostya.orlov/task-6/internal/wifi"
	wifi "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

//go:generate mockery --all --testonly --quiet --outpkg wifi_test --output

var ErrExpected = errors.New("error expected")

func TestGetAddressesSuccess(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)
	wifiServece := mywifi.New(mockWifi)

	testTable := [][]string{
		{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
		{"00:01:02:03:04:05", "aa:ab:ac:ad:ae:af"},
	}

	for _, row := range testTable {
		mockWifi.On("Interfaces").Unset()
		mockWifi.On("Interfaces").Return(mockIfaces(row), nil)

		actualAddrs, err := wifiServece.GetAddresses()

		require.Equal(t, parseMACs(row), actualAddrs)
		require.NoError(t, err)
	}
}

func TestGetAddressesError(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)
	wifiServece := mywifi.New(mockWifi)

	mockWifi.On("Interfaces").Unset()
	mockWifi.On("Interfaces").Return(nil, ErrExpected)

	actualAddrs, err := wifiServece.GetAddresses()

	require.Nil(t, actualAddrs)
	require.ErrorIs(t, err, ErrExpected)
	require.ErrorContains(t, err, "getting interfaces")
}

func TestGetNamesSuccess(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)
	wifiServece := mywifi.New(mockWifi)

	testTable := [][]string{
		{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
		{"00:01:02:03:04:05", "aa:ab:ac:ad:ae:af"},
	}

	for _, row := range testTable {
		mockWifi.On("Interfaces").Unset()
		mockWifi.On("Interfaces").Return(mockIfaces(row), nil)

		actualNames, err := wifiServece.GetNames()

		require.NoError(t, err)
		require.Equal(t, parseName(row), actualNames)
	}
}

func TestGetNamesError(t *testing.T) {
	t.Parallel()

	mockWifi := NewWiFiHandle(t)

	wifiServece := mywifi.New(mockWifi)

	mockWifi.On("Interfaces").Unset()
	mockWifi.On("Interfaces").Return(nil, ErrExpected)

	actualNames, err := wifiServece.GetNames()

	require.Nil(t, actualNames)
	require.ErrorIs(t, err, ErrExpected)
	require.ErrorContains(t, err, "getting interfaces")
}

func mockIfaces(addrs []string) []*wifi.Interface {
	interfaces := make([]*wifi.Interface, 0, len(addrs))

	for i, addr := range addrs {
		hwAddr := parseMAC(addr)
		if hwAddr == nil {
			continue
		}

		iface := &wifi.Interface{
			Index:        i + 1,
			Name:         fmt.Sprintf("eth%d", i+1),
			HardwareAddr: hwAddr,
			PHY:          1,
			Device:       1,
			Type:         wifi.InterfaceTypeAPVLAN,
			Frequency:    0,
		}

		interfaces = append(interfaces, iface)
	}

	return interfaces
}

func parseMACs(addrs []string) []net.HardwareAddr {
	hwAddrs := make([]net.HardwareAddr, 0, len(addrs))

	for _, addr := range addrs {
		hwAddrs = append(hwAddrs, parseMAC(addr))
	}

	return hwAddrs
}

func parseMAC(addr string) net.HardwareAddr {
	hwAddr, err := net.ParseMAC(addr)
	if err != nil {
		return nil
	}

	return hwAddr
}

func parseName(addrStr []string) []string {
	netNames := make([]string, 0, len(addrStr))

	for i := range addrStr {
		netNames = append(netNames, fmt.Sprintf("eth%d", i+1))
	}

	return netNames
}