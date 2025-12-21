package wifi

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
)

func TestWiFiService_GetAddresses(t *testing.T) {
	tests := []struct {
		name          string
		mockBehavior  func(m *MockWiFiHandle)
		expectedAddrs []net.HardwareAddr
		expectError   bool
	}{
		{
			name: "Success",
			mockBehavior: func(m *MockWiFiHandle) {
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
			mockBehavior: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return(nil, errors.New("driver failure"))
			},
			expectedAddrs: nil,
			expectError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockWifi := new(MockWiFiHandle)
			tc.mockBehavior(mockWifi)

			service := New(mockWifi)
			addrs, err := service.GetAddresses()

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "getting interfaces")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedAddrs, addrs)
			}
			mockWifi.AssertExpectations(t)
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	tests := []struct {
		name          string
		mockBehavior  func(m *MockWiFiHandle)
		expectedNames []string
		expectError   bool
	}{
		{
			name: "Success",
			mockBehavior: func(m *MockWiFiHandle) {
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
			mockBehavior: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return(nil, errors.New("io error"))
			},
			expectedNames: nil,
			expectError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockWifi := new(MockWiFiHandle)
			tc.mockBehavior(mockWifi)

			service := New(mockWifi)
			names, err := service.GetNames()

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "getting interfaces")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedNames, names)
			}
			mockWifi.AssertExpectations(t)
		})
	}
}
