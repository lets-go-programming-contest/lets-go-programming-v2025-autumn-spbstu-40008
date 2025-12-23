package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
	
	myWifi "github.com/task-6/internal/wifi"
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
		mockWifi.On("Interfaces").Unset()
		mockWifi.On("Interfaces").Return(mockIfaces(row.addrs), row.errExpected)

		actualAddrs, err := wifiService.GetAddresses()

		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected, "row: %d, expected error: %w, actual error: %w", i, row.errExpected, err)
			continue
		}

		require.NoError(t, err, "row: %d, error must be nil", i)
		require.Equal(t, parseMACs(row.addrs), actualAddrs, "row: %d, expected addrs: %s, actual addrs: %s", i, parseMACs(row.addrs), actualAddrs)
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
			Type:         wifi.InterfaceTypeAPVLAN,
			PHY:          1,
			Device:       1,
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