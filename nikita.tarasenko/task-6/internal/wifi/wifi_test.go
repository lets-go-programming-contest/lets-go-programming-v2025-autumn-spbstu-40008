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

	wifiModule "task-6/internal/wifi"
)

var (
	driverErr     = errors.New("driver malfunction")
	permissionErr = errors.New("insufficient permissions")
	assertionErr  = errors.New("type conversion failure")
)

type MockWiFiAdapter struct {
	mock.Mock
}

func (m *MockWiFiAdapter) NetworkInterfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	if args.Get(0) == nil {
		err := args.Error(1)
		if err != nil {
			return nil, fmt.Errorf("mock adapter error: %w", err)
		}

		return nil, nil
	}

	ifaces, conversionOk := args.Get(0).([]*wifi.Interface)
	if !conversionOk {
		err := args.Error(1)
		if err != nil {
			return nil, fmt.Errorf("interface conversion failed: %w", err)
		}

		return nil, assertionErr
	}

	err := args.Error(1)
	if err != nil {
		return ifaces, fmt.Errorf("mock adapter returned error: %w", err)
	}

	return ifaces, nil
}

func createMockInterfaces(macStrings []string) []*wifi.Interface {
	interfaces := make([]*wifi.Interface, 0, len(macStrings))

	for idx, macString := range macStrings {
		parsedAddr, parseErr := net.ParseMAC(macString)
		if parseErr != nil {
			continue
		}

		iface := &wifi.Interface{
			Index:        idx,
			Name:         "wireless" + string(rune('0'+idx)),
			HardwareAddr: parsedAddr,
		}
		interfaces = append(interfaces, iface)
	}

	return interfaces
}

func TestWiFiManager_FetchHardwareAddresses(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description string
		macStrings  []string
		mockErr     error
		shouldFail  bool
	}{
		{
			description: "successful hardware address retrieval",
			macStrings:  []string{"00:11:22:33:44:55", "AA:BB:CC:DD:EE:FF"},
			mockErr:     nil,
			shouldFail:  false,
		},
		{
			description: "no network interfaces present",
			macStrings:  []string{},
			mockErr:     nil,
			shouldFail:  false,
		},
		{
			description: "network interface error",
			macStrings:  nil,
			mockErr:     driverErr,
			shouldFail:  true,
		},
		{
			description: "malformed MAC address in list",
			macStrings:  []string{"00:11:22:33:44:55", "incorrect-format"},
			mockErr:     nil,
			shouldFail:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			mockAdapter := &MockWiFiAdapter{}

			var interfaces []*wifi.Interface
			if tc.macStrings != nil {
				interfaces = createMockInterfaces(tc.macStrings)
			}

			mockAdapter.On("NetworkInterfaces").Return(interfaces, tc.mockErr)

			manager := wifiModule.Create(mockAdapter)

			addresses, err := manager.FetchHardwareAddresses()

			if tc.shouldFail {
				require.Error(t, err)
				assert.Nil(t, addresses)
			} else {
				require.NoError(t, err)

				validCount := 0

				for _, mac := range tc.macStrings {
					if _, parseErr := net.ParseMAC(mac); parseErr == nil {
						validCount++
					}
				}

				if validCount == 0 {
					assert.Empty(t, addresses)
				} else {
					assert.Len(t, addresses, validCount)

					for _, addr := range addresses {
						assert.NotNil(t, addr)
					}
				}
			}

			mockAdapter.AssertExpectations(t)
		})
	}
}

func TestWiFiManager_FetchInterfaceNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description string
		namesList   []string
		mockErr     error
		shouldFail  bool
	}{
		{
			description: "successful interface name retrieval",
			namesList:   []string{"wireless0", "wireless1", "ethernet0"},
			mockErr:     nil,
			shouldFail:  false,
		},
		{
			description: "interface access error",
			namesList:   nil,
			mockErr:     permissionErr,
			shouldFail:  true,
		},
		{
			description: "empty interface list",
			namesList:   []string{},
			mockErr:     nil,
			shouldFail:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			mockAdapter := &MockWiFiAdapter{}

			var interfaces []*wifi.Interface

			for idx, name := range tc.namesList {
				iface := &wifi.Interface{
					Index: idx,
					Name:  name,
				}
				interfaces = append(interfaces, iface)
			}

			mockAdapter.On("NetworkInterfaces").Return(interfaces, tc.mockErr)

			manager := wifiModule.Create(mockAdapter)

			names, err := manager.FetchInterfaceNames()

			if tc.shouldFail {
				require.Error(t, err)
				assert.Nil(t, names)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.namesList, names)
			}

			mockAdapter.AssertExpectations(t)
		})
	}
}

func BenchmarkFetchHardwareAddresses(b *testing.B) {
	mockAdapter := &MockWiFiAdapter{}
	macStrings := []string{"00:11:22:33:44:55", "AA:BB:CC:DD:EE:FF"}
	interfaces := createMockInterfaces(macStrings)

	mockAdapter.On("NetworkInterfaces").Return(interfaces, nil).Times(b.N)

	manager := wifiModule.Create(mockAdapter)

	b.ResetTimer()

	for range b.N {
		_, _ = manager.FetchHardwareAddresses()
	}
}

func BenchmarkFetchInterfaceNames(b *testing.B) {
	mockAdapter := &MockWiFiAdapter{}
	namesList := []string{"wireless0", "wireless1", "ethernet0"}

	interfaces := make([]*wifi.Interface, 0, len(namesList))

	for idx, name := range namesList {
		iface := &wifi.Interface{
			Index: idx,
			Name:  name,
		}
		interfaces = append(interfaces, iface)
	}

	mockAdapter.On("NetworkInterfaces").Return(interfaces, nil).Times(b.N)

	manager := wifiModule.Create(mockAdapter)

	b.ResetTimer()

	for range b.N {
		_, _ = manager.FetchInterfaceNames()
	}
}
