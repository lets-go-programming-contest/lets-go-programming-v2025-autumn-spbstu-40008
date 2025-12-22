package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/mordw1n/task-6/internal/wifi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockWiFiHandle struct {
	mock.Mock
}

func (m *mockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	iface, ok := args.Get(0).([]*wifi.Interface)
	if !ok {
		return nil, errors.New("type assertion failed")
	}

	return iface, args.Error(1)
}

func (m *mockWiFiHandle) StationInfo(ifi *wifi.Interface) (*wifi.StationInfo, error) {
	args := m.Called(ifi)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	info, ok := args.Get(0).(*wifi.StationInfo)
	if !ok {
		return nil, errors.New("type assertion failed")
	}

	return info, args.Error(1)
}

func mockIfaces(addrs []string) []*wifi.Interface {
	interfaces := make([]*wifi.Interface, 0, len(addrs))

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
	addrs := make([]net.HardwareAddr, 0, len(macStrs))

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
	t.Parallel()

	tests := []struct {
		name          string
		setupMock     func(*mockWiFiHandle)
		expectedAddrs []string
		wantError     bool
		errorMsg      string
	}{
		{
			name: "successful get addresses",
			setupMock: func(m *mockWiFiHandle) {
				mockIfaces := mockIfaces([]string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"})
				m.On("Interfaces").Return(mockIfaces, nil)
			},
			expectedAddrs: []string{"00:11:22:33:44:55", "aa:bb:cc:dd:ee:ff"},
			wantError:     false,
		},
		{
			name: "empty interfaces",
			setupMock: func(m *mockWiFiHandle) {
				m.On("Interfaces").Return([]*wifi.Interface{}, nil)
			},
			expectedAddrs: []string{},
			wantError:     false,
		},
		{
			name: "error getting interfaces",
			setupMock: func(m *mockWiFiHandle) {
				m.On("Interfaces").Return(nil, errors.New("permission denied"))
			},
			expectedAddrs: nil,
			wantError:     true,
			errorMsg:      "failed to get interfaces",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockWiFi := &mockWiFiHandle{}
			tt.setupMock(mockWiFi)

			wifiService := wifi.New(mockWiFi)
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

			mockWiFi.AssertExpectations(t)
		})
	}
}

func TestGetInterfaceNames(t *testing.T) {
	t.Parallel()

	t.Run("get interface names", func(t *testing.T) {
		t.Parallel()

		mockWiFi := &mockWiFiHandle{}
		mockIfaces := []*wifi.Interface{
			{Name: "wlan0"},
			{Name: "wlan1"},
			{Name: "eth0"},
		}

		mockWiFi.On("Interfaces").Return(mockIfaces, nil)

		wifiService := wifi.New(mockWiFi)
		names, err := wifiService.GetInterfaceNames()

		require.NoError(t, err)
		require.Equal(t, []string{"wlan0", "wlan1", "eth0"}, names)
		mockWiFi.AssertExpectations(t)
	})

	t.Run("error getting interfaces", func(t *testing.T) {
		t.Parallel()

		mockWiFi := &mockWiFiHandle{}
		mockWiFi.On("Interfaces").Return(nil, errors.New("ioctl failed"))

		wifiService := wifi.New(mockWiFi)
		names, err := wifiService.GetInterfaceNames()

		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "failed to get interfaces")
		mockWiFi.AssertExpectations(t)
	})
}

func TestGetStationInfo(t *testing.T) {
	t.Parallel()

	t.Run("station info found", func(t *testing.T) {
		t.Parallel()

		mockWiFi := &mockWiFiHandle{}
		mockIfaces := []*wifi.Interface{
			{Name: "wlan0", Index: 1},
		}
		stationInfo := &wifi.StationInfo{}

		mockWiFi.On("Interfaces").Return(mockIfaces, nil)
		mockWiFi.On("StationInfo", mockIfaces[0]).Return(stationInfo, nil)

		wifiService := wifi.New(mockWiFi)
		info, err := wifiService.GetStationInfo("wlan0")

		require.NoError(t, err)
		require.Equal(t, stationInfo, info)
		mockWiFi.AssertExpectations(t)
	})

	t.Run("interface not found", func(t *testing.T) {
		t.Parallel()

		mockWiFi := &mockWiFiHandle{}
		mockIfaces := []*wifi.Interface{
			{Name: "wlan0", Index: 1},
		}

		mockWiFi.On("Interfaces").Return(mockIfaces, nil)

		wifiService := wifi.New(mockWiFi)
		info, err := wifiService.GetStationInfo("eth0")

		require.Error(t, err)
		require.Nil(t, info)
		require.Contains(t, err.Error(), "not found")
		mockWiFi.AssertExpectations(t)
	})
}
