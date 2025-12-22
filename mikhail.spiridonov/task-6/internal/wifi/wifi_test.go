package wifi

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"

	wifi_mocks "mordw1n/task-6/internal/mocks/wifi"
)

func mockIfaces(addrs []string) []*wifi.Interface {
	var interfaces []*wifi.Interface
	for i, addrStr := range addrs {
		hwAddr, err := net.ParseMAC(addrStr)
		if err != nil {
			continue
		}

		iface := &wifi.Interface{
			Index:        i + 1,
			Name:         fmt.Sprintf("wlan%d", i),
			HardwareAddr: hwAddr,
			PHY:          0,
			Device:       0,
			Type:         wifi.InterfaceTypeStation,
			Frequency:    2412,
		}
		interfaces = append(interfaces, iface)
	}
	return interfaces
}

func parseMACs(macStrs []string) []net.HardwareAddr {
	var addrs []net.HardwareAddr
	for _, macStr := range macStrs {
		hwAddr, err := net.ParseMAC(macStr)
		if err != nil {
			continue
		}
		addrs = append(addrs, hwAddr)
	}
	return addrs
}

func TestGetAddresses(t *testing.T) {
	mockWiFi := wifi_mocks.NewWiFiHandle(t)
	wifiService := New(mockWiFi)

	tests := []struct {
		name          string
		setupMock     func()
		expectedAddrs []string
		wantError     bool
		errorMsg      string
	}{
		{
			name: "successful get addresses",
			setupMock: func() {
				mockIfaces := mockIfaces([]string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"})
				mockWiFi.EXPECT().Interfaces().Return(mockIfaces, nil)
			},
			expectedAddrs: []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
			wantError:     false,
		},
		{
			name: "empty interfaces",
			setupMock: func() {
				mockWiFi.EXPECT().Interfaces().Return([]*wifi.Interface{}, nil)
			},
			expectedAddrs: []string{},
			wantError:     false,
		},
		{
			name: "getting interfaces",
			setupMock: func() {
				mockWiFi.EXPECT().Interfaces().
					Return(nil, errors.New("permission denied"))
			},
			expectedAddrs: nil,
			wantError:     true,
			errorMsg:      "get interfaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			addrs, err := wifiService.GetAddresses()

			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
				require.Nil(t, addrs)
			} else {
				require.NoError(t, err)
				expectedMACs := parseMACs(tt.expectedAddrs)
				require.Equal(t, expectedMACs, addrs)
			}
		})
	}
}

func TestGetInterfaceNames(t *testing.T) {
	mockWiFi := wifi_mocks.NewWiFiHandle(t)
	wifiService := New(mockWiFi)

	t.Run("get interface names", func(t *testing.T) {
		mockIfaces := []*wifi.Interface{
			{Name: "wlan0"},
			{Name: "wlan1"},
			{Name: "eth0"},
		}
		mockWiFi.EXPECT().Interfaces().Return(mockIfaces, nil)

		names, err := wifiService.GetInterfaceNames()

		require.NoError(t, err)
		require.Equal(t, []string{"wlan0", "wlan1", "eth0"}, names)
	})

	t.Run("getting interfaces", func(t *testing.T) {
		mockWiFi.EXPECT().Interfaces().
			Return(nil, errors.New("ioctl"))

		names, err := wifiService.GetInterfaceNames()

		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "get interfaces")
	})
}

func TestGetStationInfo(t *testing.T) {
	mockWiFi := wifi_mocks.NewWiFiHandle(t)
	wifiService := New(mockWiFi)

	t.Run("station info found", func(t *testing.T) {
		mockIfaces := []*wifi.Interface{
			{Name: "wlan0", Index: 1},
		}
		stationInfo := &wifi.StationInfo{}

		mockWiFi.EXPECT().Interfaces().Return(mockIfaces, nil)
		mockWiFi.EXPECT().StationInfo(mockIfaces[0]).Return(stationInfo, nil)

		info, err := wifiService.GetStationInfo("wlan0")

		require.NoError(t, err)
		require.Equal(t, stationInfo, info)
	})

	t.Run("interface not found", func(t *testing.T) {
		mockIfaces := []*wifi.Interface{
			{Name: "wlan0", Index: 1},
		}

		mockWiFi.EXPECT().Interfaces().Return(mockIfaces, nil)

		info, err := wifiService.GetStationInfo("eth0")

		require.Error(t, err)
		require.Nil(t, info)
		require.Contains(t, err.Error(), "not found")
	})
}
