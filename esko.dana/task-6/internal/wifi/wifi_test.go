package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	wifiPkg "esko.dana/task-6/internal/wifi"
)

var (
	errFailedInterfaces = errors.New("failed to get interfaces")
	errPermissionDenied = errors.New("permission denied")
)

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	createInterfaces := func(macs []string) []*wifi.Interface {
		interfaces := make([]*wifi.Interface, 0, len(macs))

		for i, macStr := range macs {
			mac, _ := net.ParseMAC(macStr)

			interfaces = append(interfaces, &wifi.Interface{
				Index:        i,
				Name:         "wlan",
				HardwareAddr: mac,
			})
		}

		return interfaces
	}

	testCases := []struct {
		name        string
		setupMock   func(*MockWiFiHandle)
		expected    []net.HardwareAddr
		expectedErr string
	}{
		{
			name: "success",
			setupMock: func(m *MockWiFiHandle) {
				m.On("Interfaces").
					Return(createInterfaces([]string{"00:11:22:33:44:55"}), nil).
					Once()
			},
			expected: func() []net.HardwareAddr {
				mac, _ := net.ParseMAC("00:11:22:33:44:55")
				return []net.HardwareAddr{mac}
			}(),
		},
		{
			name: "error",
			setupMock: func(m *MockWiFiHandle) {
				m.On("Interfaces").
					Return(nil, errFailedInterfaces).
					Once()
			},
			expectedErr: "getting interfaces:",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockWiFi := &MockWiFiHandle{}
			service := wifiPkg.New(mockWiFi)

			tc.setupMock(mockWiFi)

			result, err := service.GetAddresses()

			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}

			mockWiFi.AssertExpectations(t)
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupMock   func(*MockWiFiHandle)
		expected    []string
		expectedErr string
	}{
		{
			name: "success",
			setupMock: func(m *MockWiFiHandle) {
				m.On("Interfaces").
					Return([]*wifi.Interface{
						{Name: "wlan0"},
						{Name: "eth0"},
					}, nil).
					Once()
			},
			expected: []string{"wlan0", "eth0"},
		},
		{
			name: "error",
			setupMock: func(m *MockWiFiHandle) {
				m.On("Interfaces").
					Return(nil, errPermissionDenied).
					Once()
			},
			expectedErr: "getting interfaces:",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockWiFi := &MockWiFiHandle{}
			service := wifiPkg.New(mockWiFi)

			tc.setupMock(mockWiFi)

			result, err := service.GetNames()

			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}

			mockWiFi.AssertExpectations(t)
		})
	}
}
