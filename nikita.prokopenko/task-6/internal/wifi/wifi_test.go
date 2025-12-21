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

	iwifi "github.com/Czeeen/lets-go-programming-v2025-autumn-spbstu-40008/nikita.prokopenko/task-6/internal/wifi"
)

var (
	errInterfaceError = errors.New("interface access error")
	errPermission     = errors.New("permission denied")
	errTypeAssertion  = errors.New("type assertion failed for interface slice")
)

type MockInterfaceSource struct {
	mock.Mock
}

func (m *MockInterfaceSource) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	if args.Get(0) == nil {
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("mock error: %w", err)
		}
		return nil, nil
	}
	interfaces, ok := args.Get(0).([]*wifi.Interface)
	if !ok {
		return nil, errTypeAssertion
	}
	if err := args.Error(1); err != nil {
		return interfaces, fmt.Errorf("mock error: %w", err)
	}
	return interfaces, nil
}

func createTestInterface(name, macStr string) *wifi.Interface {
	mac, _ := net.ParseMAC(macStr)
	return &wifi.Interface{
		Name:         name,
		HardwareAddr: mac,
	}
}

func TestNetworkManager_GetMACAddresses(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name           string
		mockSetup      func(*MockInterfaceSource)
		expectedMACs   []string
		expectError    bool
		errorSubstring string
	}{
		{
			name: "valid MAC addresses",
			mockSetup: func(m *MockInterfaceSource) {
				interfaces := []*wifi.Interface{
					createTestInterface("eth0", "01:23:45:67:89:ab"),
					createTestInterface("wlan0", "cd:ef:01:23:45:67"),
				}
				m.On("Interfaces").Return(interfaces, nil).Once()
			},
			expectedMACs: []string{"01:23:45:67:89:ab", "cd:ef:01:23:45:67"},
		},
		{
			name: "interface fetch error",
			mockSetup: func(m *MockInterfaceSource) {
				m.On("Interfaces").Return(nil, errInterfaceError).Once()
			},
			expectError:    true,
			errorSubstring: "failed to fetch interfaces",
		},
		{
			name: "no valid MACs",
			mockSetup: func(m *MockInterfaceSource) {
				interfaces := []*wifi.Interface{
					{Name: "lo", HardwareAddr: net.HardwareAddr{}},
					{Name: "dummy", HardwareAddr: net.HardwareAddr{0x00}},
				}
				m.On("Interfaces").Return(interfaces, nil).Once()
			},
			expectError:    true,
			errorSubstring: "no valid MAC addresses",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockSource := new(MockInterfaceSource)
			manager := iwifi.CreateManager(mockSource)
			tc.mockSetup(mockSource)
			macs, err := manager.GetMACAddresses()
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorSubstring)
				assert.Nil(t, macs)
			} else {
				require.NoError(t, err)
				require.Len(t, macs, len(tc.expectedMACs))
				for i, expected := range tc.expectedMACs {
					assert.Equal(t, expected, macs[i].String())
				}
			}
			mockSource.AssertExpectations(t)
		})
	}
}

func TestNetworkManager_GetInterfaceNames(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name           string
		mockSetup      func(*MockInterfaceSource)
		expectedNames  []string
		expectError    bool
		errorSubstring string
	}{
		{
			name: "valid interface names",
			mockSetup: func(m *MockInterfaceSource) {
				interfaces := []*wifi.Interface{
					{Name: "eth0"},
					{Name: "wlan1"},
					{Name: "docker0"},
				}
				m.On("Interfaces").Return(interfaces, nil).Once()
			},
			expectedNames: []string{"eth0", "wlan1", "docker0"},
		},
		{
			name: "permission error",
			mockSetup: func(m *MockInterfaceSource) {
				m.On("Interfaces").Return(nil, errPermission).Once()
			},
			expectError:    true,
			errorSubstring: "failed to fetch interfaces",
		},
		{
			name: "all names empty",
			mockSetup: func(m *MockInterfaceSource) {
				interfaces := []*wifi.Interface{
					{Name: ""},
					{Name: " "},
				}
				m.On("Interfaces").Return(interfaces, nil).Once()
			},
			expectError:    true,
			errorSubstring: "all names empty",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockSource := new(MockInterfaceSource)
			manager := iwifi.CreateManager(mockSource)
			tc.mockSetup(mockSource)
			names, err := manager.GetInterfaceNames()
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorSubstring)
				assert.Nil(t, names)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedNames, names)
			}
			mockSource.AssertExpectations(t)
		})
	}
}
