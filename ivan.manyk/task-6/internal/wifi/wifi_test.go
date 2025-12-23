package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	wifiPkg "task-6/internal/wifi"
)

var (
	errDriver            = errors.New("driver error")
	errPermission        = errors.New("permission denied")
	errTypeAssertion     = errors.New("type assertion failed")
)

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	if args.Get(0) == nil {
		if args.Error(1) != nil {
			return nil, fmt.Errorf("mock error: %w", args.Error(1))
		}

		return nil, args.Error(1)
	}

	ifaces, ok := args.Get(0).([]*wifi.Interface)
	if !ok {
		if args.Error(1) != nil {
			return nil, fmt.Errorf("type assertion failed: %w", args.Error(1))
		}

		return nil, errTypeAssertion
	}

	return ifaces, args.Error(1)
}

func mockIfaces(macAddrs []string) []*wifi.Interface {
	interfaces := make([]*wifi.Interface, 0, len(macAddrs))

	for i, mac := range macAddrs {
		hwAddr, err := net.ParseMAC(mac)
		if err != nil {
			continue
		}

		iface := &wifi.Interface{
			Index:        i,
			Name:         "wlan" + string(rune('0'+i)),
			HardwareAddr: hwAddr,
		}
		interfaces = append(interfaces, iface)
	}

	return interfaces
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		macAddrs    []string
		mockError   error
		expectError bool
	}{
		{
			name:        "successful get MAC addresses",
			macAddrs:    []string{"00:11:22:33:44:55", "AA:BB:CC:DD:EE:FF"},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "empty interface list",
			macAddrs:    []string{},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "interface error",
			macAddrs:    nil,
			mockError:   errDriver,
			expectError: true,
		},
		{
			name:        "invalid MAC address",
			macAddrs:    []string{"00:11:22:33:44:55", "invalid-mac"},
			mockError:   nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockWiFi := &MockWiFiHandle{}

			var interfaces []*wifi.Interface
			if tt.macAddrs != nil {
				interfaces = mockIfaces(tt.macAddrs)
			}

			mockWiFi.On("Interfaces").Return(interfaces, tt.mockError)

			service := wifiPkg.New(mockWiFi)

			addrs, err := service.GetAddresses()

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, addrs)
			} else {
				require.NoError(t, err)

				expectedCount := 0

				for _, mac := range tt.macAddrs {
					if _, err := net.ParseMAC(mac); err == nil {
						expectedCount++
					}
				}

				if expectedCount == 0 {
					assert.Empty(t, addrs)
				} else {
					assert.Len(t, addrs, expectedCount)

					for _, addr := range addrs {
						assert.NotNil(t, addr)
					}
				}
			}

			mockWiFi.AssertExpectations(t)
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		ifaceNames  []string
		mockError   error
		expectError bool
	}{
		{
			name:        "successful get interface names",
			ifaceNames:  []string{"wlan0", "wlan1", "eth0"},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "interface error",
			ifaceNames:  nil,
			mockError:   errPermission,
			expectError: true,
		},
		{
			name:        "empty interface list",
			ifaceNames:  []string{},
			mockError:   nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockWiFi := &MockWiFiHandle{}

			var interfaces []*wifi.Interface

			for i, name := range tt.ifaceNames {
				iface := &wifi.Interface{
					Index: i,
					Name:  name,
				}
				interfaces = append(interfaces, iface)
			}

			mockWiFi.On("Interfaces").Return(interfaces, tt.mockError)

			service := wifiPkg.New(mockWiFi)

			names, err := service.GetNames()

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, names)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.ifaceNames, names)
			}

			mockWiFi.AssertExpectations(t)
		})
	}
}

func BenchmarkGetAddresses(b *testing.B) {
	mockWiFi := &MockWiFiHandle{}
	macAddrs := []string{"00:11:22:33:44:55", "AA:BB:CC:DD:EE:FF"}
	interfaces := mockIfaces(macAddrs)

	mockWiFi.On("Interfaces").Return(interfaces, nil).Times(b.N)

	service := wifiPkg.New(mockWiFi)

	b.ResetTimer()

	for range b.N {
		_, _ = service.GetAddresses()
	}
}

func BenchmarkGetNames(b *testing.B) {
	mockWiFi := &MockWiFiHandle{}
	ifaceNames := []string{"wlan0", "wlan1", "eth0"}

	interfaces := make([]*wifi.Interface, 0, len(ifaceNames))

	for i, name := range ifaceNames {
		iface := &wifi.Interface{
			Index: i,
			Name:  name,
		}
		interfaces = append(interfaces, iface)
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil).Times(b.N)

	service := wifiPkg.New(mockWiFi)

	b.ResetTimer()

	for range b.N {
		_, _ = service.GetNames()
	}
}
