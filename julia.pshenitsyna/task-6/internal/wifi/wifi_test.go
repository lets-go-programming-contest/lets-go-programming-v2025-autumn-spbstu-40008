package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	wifiPkg "github.com/julia.pshenitsyna/task-6/internal/wifi"
)

var (
	errFailedInterfaces = errors.New("failed to get interfaces")
	errPermissionDenied = errors.New("permission denied")
)

func TestWiFiService_GetAddresses(t *testing.T) {
	mockWiFi := &MockWiFiHandle{}

	createMockInterfaces := func(macAddresses []string) []*wifi.Interface {
		interfaces := make([]*wifi.Interface, 0, len(macAddresses))
		for i, macStr := range macAddresses {
			mac, _ := net.ParseMAC(macStr)
			interfaces = append(interfaces, &wifi.Interface{
				Index:        i,
				Name:         "wlan" + string(rune('0'+i)),
				HardwareAddr: mac,
			})
		}
		return interfaces
	}

	testCases := []struct {
		name        string
		setupMock   func()
		expected    []net.HardwareAddr
		expectedErr error
	}{
		{
			name: "success - return MAC",
			setupMock: func() {
				interfaces := createMockInterfaces([]string{
					"00:11:22:33:44:55",
					"aa:bb:cc:dd:ee:ff",
				})
				mockWiFi.On("Interfaces").Return(interfaces, nil).Once()
			},
			expected: func() []net.HardwareAddr {
				mac1, _ := net.ParseMAC("00:11:22:33:44:55")
				mac2, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")
				return []net.HardwareAddr{mac1, mac2}
			}(),
		},
		{
			name: "error - faild getting interface",
			setupMock: func() {
				mockWiFi.On("Interfaces").Return(nil, errFailedInterfaces).Once()
			},
			expectedErr: errors.New("getting interfaces: " + errFailedInterfaces.Error()),
		},
		{
			name: "success - empty interfaces list",
			setupMock: func() {
				mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, nil).Once()
			},
			expected: []net.HardwareAddr{},
		},
		{
			name: "success - interface without MAC",
			setupMock: func() {
				interfaces := []*wifi.Interface{
					{
						Index:        0,
						Name:         "wlan0",
						HardwareAddr: nil,
					},
				}
				mockWiFi.On("Interfaces").Return(interfaces, nil).Once()
			},
			expected: []net.HardwareAddr{nil},
		},
		{
			name: "success - single interface with MAC",
			setupMock: func() {
				mac, _ := net.ParseMAC("00:11:22:33:44:55")
				interfaces := []*wifi.Interface{
					{
						Index:        0,
						Name:         "wlan0",
						HardwareAddr: mac,
					},
				}
				mockWiFi.On("Interfaces").Return(interfaces, nil).Once()
			},
			expected: func() []net.HardwareAddr {
				mac, _ := net.ParseMAC("00:11:22:33:44:55")
				return []net.HardwareAddr{mac}
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockWiFi := NewMockWiFiHandle(t)
			wifiService := wifiPkg.New(mockWiFi)
			
			tc.setupMock()

			result, err := wifiService.GetAddresses()

			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr.Error())
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
	mockWiFi := &MockWiFiHandle{}

	createMockInterfaces := func(names []string) []*wifi.Interface {
		interfaces := make([]*wifi.Interface, 0, len(names))
		for i, name := range names {
			interfaces = append(interfaces, &wifi.Interface{
				Index: i,
				Name:  name,
			})
		}
		return interfaces
	}

	testCases := []struct {
		name        string
		setupMock   func()
		expected    []string
		expectedErr error
	}{
		{
			name: "success - return interfase's names",
			setupMock: func() {
				interfaces := createMockInterfaces([]string{
					"wlan0",
					"wlan1",
					"eth0",
				})
				mockWiFi.On("Interfaces").Return(interfaces, nil).Once()
			},
			expected: []string{"wlan0", "wlan1", "eth0"},
		},
		{
			name: "error - getting interface error",
			setupMock: func() {
				mockWiFi.On("Interfaces").Return(nil, errPermissionDenied).Once()
			},
			expectedErr: errors.New("getting interfaces: " + errPermissionDenied.Error()),
		},
		{
			name: "success - empty interfaces list",
			setupMock: func() {
				mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, nil).Once()
			},
			expected: []string{},
		},
		{
			name: "success - same names of interfaces",
			setupMock: func() {
				interfaces := createMockInterfaces([]string{
					"wlan0",
					"wlan0",
					"eth0",
				})
				mockWiFi.On("Interfaces").Return(interfaces, nil).Once()
			},
			expected: []string{"wlan0", "wlan0", "eth0"},
		},
		{
			name: "success - single interface name",
			setupMock: func() {
				interfaces := []*wifi.Interface{
					{
						Index: 0,
						Name:  "wlan0",
					},
				}
				mockWiFi.On("Interfaces").Return(interfaces, nil).Once()
			},
			expected: []string{"wlan0"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			
			mockWiFi := NewMockWiFiHandle(t)
			wifiService := wifiPkg.New(mockWiFi)

			tc.setupMock()

			result, err := wifiService.GetNames()

			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr.Error())
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}

			mockWiFi.AssertExpectations(t)
		})
	}
}

func TestNewWiFiService(t *testing.T) {
	t.Parallel()

	mockWiFi := &MockWiFiHandle{}
	service := wifiPkg.New(mockWiFi)

	assert.NotNil(t, service)
	assert.Equal(t, mockWiFi, service.WiFi)
}
