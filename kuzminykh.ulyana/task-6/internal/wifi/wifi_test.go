package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"

	wifiExt "github.com/kuzminykh.ulyana/task-6/internal/wifi"
)

type rowTestSysInfo struct {
	addrs       []string
	errExpected error
}

var testTable = []rowTestSysInfo{
	{
		addrs: []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
	},
	{
		errExpected: errors.New("failed to get WiFi interfaces"),
	},
}

func TestGetAddresses(t *testing.T) {
	mockWifi := NewWiFiHandle(t)
	service := wifiExt.New(mockWifi)

	for i, row := range testTable {
		mockWifi.On("Interfaces").Unset()
		mockWifi.On("Interfaces").Return(mockIfaces(row.addrs), row.errExpected)

		addrs, err := service.GetAddresses()

		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected, "row: %d", i)
			continue
		}

		require.NoError(t, err, "row: %d", i)
		require.Equal(t, parseMACs(row.addrs), addrs, "row: %d", i)
	}
}

func TestGetNames(t *testing.T) {
	mockWifi := NewWiFiHandle(t)
	service := wifiExt.New(mockWifi)

	for i, row := range testTable {
		mockWifi.On("Interfaces").Unset()
		mockWifi.On("Interfaces").Return(mockIfaces(row.addrs), row.errExpected)

		names, err := service.GetNames()

		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected, "row: %d", i)
			continue
		}

		expectedNames := make([]string, len(row.addrs))
		for j := range row.addrs {
			expectedNames[j] = fmt.Sprintf("eth%d", j+1)
		}

		require.NoError(t, err, "row: %d", i)
		require.Equal(t, expectedNames, names, "row: %d", i)
	}
}

func mockIfaces(addrs []string) []*wifi.Interface {
	var interfaces []*wifi.Interface
	for i, addrStr := range addrs {
		hwAddr := parseMAC(addrStr)
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

func parseMACs(macStr []string) []net.HardwareAddr {
	var addrs []net.HardwareAddr
	for _, addr := range macStr {
		addrs = append(addrs, parseMAC(addr))
	}
	return addrs
}

func parseMAC(macStr string) net.HardwareAddr {
	hwAddr, err := net.ParseMAC(macStr)
	if err != nil {
		return nil
	}
	return hwAddr
}
