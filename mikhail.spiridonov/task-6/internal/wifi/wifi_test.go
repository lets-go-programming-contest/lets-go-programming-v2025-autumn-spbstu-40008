package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/mordw1n/task-6/internal/wifi"
	"github.com/stretchr/testify/require"
)

type MockWiFiHandle struct {
	interfaces []*wifi.Interface
	err        error
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	return m.interfaces, m.err
}

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		interfaces  []*wifi.Interface
		mockErr     error
		wantAddrs   []net.HardwareAddr
		wantError   bool
		errorMsg    string
	}{
		{
			name: "successful get addresses",
			interfaces: []*wifi.Interface{
				{
					Name:         "wlan0",
					HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
				},
				{
					Name:         "wlan1",
					HardwareAddr: net.HardwareAddr{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF},
				},
			},
			wantAddrs: []net.HardwareAddr{
				{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
				{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF},
			},
			wantError: false,
		},
		{
			name:       "empty interfaces",
			interfaces: []*wifi.Interface{},
			wantAddrs:  []net.HardwareAddr{},
			wantError:  false,
		},
		{
			name: "interface with nil hardware address",
			interfaces: []*wifi.Interface{
				{
					Name:         "wlan0",
					HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
				},
				{
					Name:         "wlan1",
					HardwareAddr: nil,
				},
			},
			wantAddrs: []net.HardwareAddr{
				{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
				nil,
			},
			wantError: false,
		},
		{
			name:      "error getting interfaces",
			mockErr:   errors.New("interface error"),
			wantAddrs: nil,
			wantError: true,
			errorMsg:  "getting interfaces:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockHandle := &MockWiFiHandle{
				interfaces: tt.interfaces,
				err:        tt.mockErr,
			}

			service := wifi.New(mockHandle)
			addrs, err := service.GetAddresses()

			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantAddrs, addrs)
			}
		})
	}
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		interfaces []*wifi.Interface
		mockErr    error
		wantNames  []string
		wantError  bool
		errorMsg   string
	}{
		{
			name: "successful get names",
			interfaces: []*wifi.Interface{
				{Name: "wlan0"},
				{Name: "wlan1"},
				{Name: "wlp2s0"},
			},
			wantNames: []string{"wlan0", "wlan1", "wlp2s0"},
			wantError: false,
		},
		{
			name:       "empty interfaces",
			interfaces: []*wifi.Interface{},
			wantNames:  []string{},
			wantError:  false,
		},
		{
			name:      "error getting interfaces",
			mockErr:   errors.New("interface error"),
			wantNames: nil,
			wantError: true,
			errorMsg:  "getting interfaces:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockHandle := &MockWiFiHandle{
				interfaces: tt.interfaces,
				err:        tt.mockErr,
			}

			service := wifi.New(mockHandle)
			names, err := service.GetNames()

			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantNames, names)
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockHandle := &MockWiFiHandle{}
	service := wifi.New(mockHandle)
	require.NotNil(t, service)
}
