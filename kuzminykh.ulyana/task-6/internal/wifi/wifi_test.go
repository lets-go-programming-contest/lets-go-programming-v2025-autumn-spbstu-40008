package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	wifiExt "github.com/kuzminykh.ulyana/task-6/internal/wifi"
)

var errWiFi = errors.New("failed to get WiFi interfaces")

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupMock   func(*WiFiHandle)
		expected    []net.HardwareAddr
		expectedErr string
	}{
		{
			name: "success - return addresses",
			setupMock: func(mock *WiFiHandle) {
				addr1, _ := net.ParseMAC("00:11:22:33:44:55")
				addr2, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")
				mock.On("Interfaces").Return([]*wifi.Interface{
					{HardwareAddr: addr1},
					{HardwareAddr: addr2},
				}, nil)
			},
			expected: []net.HardwareAddr{
				mustParseMAC("00:11:22:33:44:55"),
				mustParseMAC("aa:bb:cc:dd:ee:ff"),
			},
		},
		{
			name: "error - wifi interfaces error",
			setupMock: func(mock *WiFiHandle) {
				mock.On("Interfaces").Return([]*wifi.Interface(nil), errWiFi)
			},
			expectedErr: "getting interfaces",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHandle := NewWiFiHandle(t)
			tc.setupMock(mockHandle)

			service := wifiExt.New(mockHandle)
			result, err := service.GetAddresses()

			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}

			mockHandle.AssertExpectations(t)
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupMock   func(*WiFiHandle)
		expected    []string
		expectedErr string
	}{
		{
			name: "success - return names",
			setupMock: func(mock *WiFiHandle) {
				addr1, _ := net.ParseMAC("00:11:22:33:44:55")
				mock.On("Interfaces").Return([]*wifi.Interface{
					{Name: "wlan0", HardwareAddr: addr1},
					{Name: "wlan1", HardwareAddr: addr1},
				}, nil)
			},
			expected: []string{"wlan0", "wlan1"},
		},
		{
			name: "error - wifi interfaces error",
			setupMock: func(mock *WiFiHandle) {
				mock.On("Interfaces").Return([]*wifi.Interface(nil), errWiFi)
			},
			expectedErr: "getting interfaces",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHandle := NewWiFiHandle(t)
			tc.setupMock(mockHandle)

			service := wifiExt.New(mockHandle)
			result, err := service.GetNames()

			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}

			mockHandle.AssertExpectations(t)
		})
	}
}

func mustParseMAC(s string) net.HardwareAddr {
	addr, err := net.ParseMAC(s)
	if err != nil {
		panic(err)
	}
	return addr
}
