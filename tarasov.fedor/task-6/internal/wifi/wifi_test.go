package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mywifi "github.com/task-6/internal/wifi"
)

var errMock = errors.New("mock error")

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mockBehavior  func(m *mywifi.MockWiFiHandle)
		expectedAddrs []net.HardwareAddr
		expectError   bool
	}{
		{
			name: "Success",
			mockBehavior: func(m *mywifi.MockWiFiHandle) {
				hwAddr, _ := net.ParseMAC("00:00:5e:00:53:01")
				ifaces := []*wifi.Interface{
					{HardwareAddr: hwAddr},
				}
				m.On("Interfaces").Return(ifaces, nil)
			},
			expectedAddrs: []net.HardwareAddr{
				{0x00, 0x00, 0x5e, 0x00, 0x53, 0x01},
			},
			expectError: false,
		},
		{
			name: "Error Getting Interfaces",
			mockBehavior: func(m *mywifi.MockWiFiHandle) {
				m.On("Interfaces").Return(nil, errMock)
			},
			expectedAddrs: nil,
			expectError:   true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockWifi := new(mywifi.MockWiFiHandle)
			tc.mockBehavior(mockWifi)

			service := mywifi.New(mockWifi)
			addrs, err := service.GetAddresses()

			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "getting interfaces")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedAddrs, addrs)
			}

			mockWifi.AssertExpectations(t)
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mockBehavior  func(m *mywifi.MockWiFiHandle)
		expectedNames []string
		expectError   bool
	}{
		{
			name: "Success",
			mockBehavior: func(m *mywifi.MockWiFiHandle) {
				ifaces := []*wifi.Interface{
					{Name: "wlan0"},
					{Name: "wlan1"},
				}
				m.On("Interfaces").Return(ifaces, nil)
			},
			expectedNames: []string{"wlan0", "wlan1"},
			expectError:   false,
		},
		{
			name: "Error Getting Interfaces",
			mockBehavior: func(m *mywifi.MockWiFiHandle) {
				m.On("Interfaces").Return(nil, errMock)
			},
			expectedNames: nil,
			expectError:   true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockWifi := new(mywifi.MockWiFiHandle)
			tc.mockBehavior(mockWifi)

			service := mywifi.New(mockWifi)
			names, err := service.GetNames()

			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "getting interfaces")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedNames, names)
			}

			mockWifi.AssertExpectations(t)
		})
	}
}
