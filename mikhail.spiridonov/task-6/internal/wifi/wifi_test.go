package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	mywifi "github.com/mordw1n/task-6/internal/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errMock = errors.New("mock error")

type MockWiFiHandle struct {
	mock.Mock
}

func (_m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	ret := _m.Called()

	var r0 []*wifi.Interface
	if rf, ok := ret.Get(0).(func() []*wifi.Interface); ok {
		r0 = rf()
	} else if ret.Get(0) != nil {
		if val, ok := ret.Get(0).([]*wifi.Interface); ok {
			r0 = val
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

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
			name: "Success with multiple interfaces",
			mockBehavior: func(m *MockWiFiHandle) {
				hwAddr1, _ := net.ParseMAC("00:11:22:33:44:55")
				hwAddr2, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")
				ifaces := []*wifi.Interface{
					{HardwareAddr: hwAddr1},
					{HardwareAddr: hwAddr2},
				}
				m.On("Interfaces").Return(ifaces, nil)
			},
			expectedAddrs: []net.HardwareAddr{
				{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
				{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
			},
			expectError: false,
		},
		{
			name: "Interface with nil hardware address",
			mockBehavior: func(m *MockWiFiHandle) {
				hwAddr, _ := net.ParseMAC("00:11:22:33:44:55")
				ifaces := []*wifi.Interface{
					{HardwareAddr: hwAddr},
					{HardwareAddr: nil},
				}
				m.On("Interfaces").Return(ifaces, nil)
			},
			expectedAddrs: []net.HardwareAddr{
				{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
				nil,
			},
			expectError: false,
		},
		{
			name: "Error Getting Interfaces",
			mockBehavior: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return(nil, errMock)
			},
			expectedAddrs: nil,
			expectError:   true,
		},
		{
			name: "Empty interfaces list",
			mockBehavior: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return([]*wifi.Interface{}, nil)
			},
			expectedAddrs: []net.HardwareAddr{},
			expectError:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockWifi := new(MockWiFiHandle)
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
			name: "Success with multiple interfaces",
			mockBehavior: func(m *MockWiFiHandle) {
				ifaces := []*wifi.Interface{
					{Name: "wlan0"},
					{Name: "wlan1"},
					{Name: "wlp2s0"},
					{Name: "wlp3s0"},
				}
				m.On("Interfaces").Return(ifaces, nil)
			},
			expectedNames: []string{"wlan0", "wlan1", "wlp2s0", "wlp3s0"},
			expectError:   false,
		},
		{
			name: "Error Getting Interfaces",
			mockBehavior: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return(nil, errMock)
			},
			expectedNames: nil,
			expectError:   true,
		},
		{
			name: "Empty interfaces list",
			mockBehavior: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return([]*wifi.Interface{}, nil)
			},
			expectedNames: []string{},
			expectError:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockWifi := new(MockWiFiHandle)
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

func TestNew(t *testing.T) {
	t.Parallel()

	mockWifi := new(MockWiFiHandle)
	service := mywifi.New(mockWifi)
	require.NotNil(t, service)
}
