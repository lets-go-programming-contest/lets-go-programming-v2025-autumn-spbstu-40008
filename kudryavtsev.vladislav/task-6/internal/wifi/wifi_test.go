package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	myWifi "example_mock/internal/wifi"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
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
		errExpected: errors.New("ExpectedError"),
	},
}

func TestGetAddresses(t *testing.T) {
	mockWifi := NewWiFi(t)
	wifiService := myWifi.New(mockWifi)

	for i, row := range testTable {
		mockWifi.On("Interfaces").Return(mockIfaces(row.addrs), row.errExpected).Once()

		actualAddrs, err := wifiService.GetAddresses()

		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected, "row: %d", i)
			continue
		}

		require.NoError(t, err, "row: %d, error must be nil", i)

		expectedAddrs := parseMACs(row.addrs)
		require.Equal(t, expectedAddrs, actualAddrs, "row: %d", i)
	}
}

func mockIfaces(addrs []string) []*wifi.Interface {
	if addrs == nil {
		return nil
	}
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
			Type:         wifi.InterfaceTypeAPVLAN,
			PHY:          1,
			Device:       1,
			Frequency:    0,
		}
		interfaces = append(interfaces, iface)
	}
	return interfaces
}

func parseMACs(macStrs []string) []net.HardwareAddr {
	var addrs []net.HardwareAddr
	for _, addr := range macStrs {
		if res := parseMAC(addr); res != nil {
			addrs = append(addrs, res)
		}
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
