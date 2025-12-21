package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	wifiPkg "github.com/Czeeen/lets-go-programming-v2025-autumn-spbstu-40008/nikita.prokopenko/task-6/internal/wifi"
)

var (
	errInterfaceError = errors.New("interface access error")
	errPermission     = errors.New("permission denied")
)

type MockProvider struct {
	interfaces []*wifi.Interface
	err        error
}

func (m *MockProvider) Interfaces() ([]*wifi.Interface, error) {
	return m.interfaces, m.err
}

func TestNetworkService_GetAddresses(t *testing.T) {
	t.Parallel()

	createTestInterface := func(name, macAddress string) *wifi.Interface {
		mac, _ := net.ParseMAC(macAddress)

		return &wifi.Interface{
			Name:         name,
			HardwareAddr: mac,
		}
	}

	cases := []struct {
		name           string
		provider       *MockProvider
		expected       []string
		expectError    bool
		errorSubstring string
	}{
		{
			name: "success with valid interfaces",
			provider: &MockProvider{
				interfaces: []*wifi.Interface{
					createTestInterface("eth0", "01:02:03:04:05:06"),
					createTestInterface("wlan0", "01:02:03:04:05:07"),
				},
			},
			expected: []string{"01:02:03:04:05:06", "01:02:03:04:05:07"},
		},
		{
			name: "error fetching interfaces",
			provider: &MockProvider{
				err: errInterfaceError,
			},
			expectError:    true,
			errorSubstring: "failed to fetch interfaces",
		},
		{
			name: "no valid interfaces",
			provider: &MockProvider{
				interfaces: []*wifi.Interface{},
			},
			expectError:    true,
			errorSubstring: "no valid network interfaces found",
		},
		{
			name: "no valid MAC addresses",
			provider: &MockProvider{
				interfaces: []*wifi.Interface{
					{Name: "lo", HardwareAddr: net.HardwareAddr{}},
				},
			},
			expectError:    true,
			errorSubstring: "no valid network interfaces found",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			service := wifiPkg.New(tc.provider)
			result, err := service.GetAddresses()

			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorSubstring)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Len(t, result, len(tc.expected))

				for i, expected := range tc.expected {
					assert.Equal(t, expected, result[i].String())
				}
			}
		})
	}
}

func TestNetworkService_GetNames(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		provider       *MockProvider
		expected       []string
		expectError    bool
		errorSubstring string
	}{
		{
			name: "success with valid interfaces",
			provider: &MockProvider{
				interfaces: []*wifi.Interface{
					{Name: "eth0"},
					{Name: "wlan0"},
				},
			},
			expected: []string{"eth0", "wlan0"},
		},
		{
			name: "error fetching interfaces",
			provider: &MockProvider{
				err: errPermission,
			},
			expectError:    true,
			errorSubstring: "failed to fetch interfaces",
		},
		{
			name: "no valid interfaces",
			provider: &MockProvider{
				interfaces: []*wifi.Interface{},
			},
			expectError:    true,
			errorSubstring: "no valid network interfaces found",
		},
		{
			name: "no valid names",
			provider: &MockProvider{
				interfaces: []*wifi.Interface{
					{Name: ""},
				},
			},
			expectError:    true,
			errorSubstring: "no valid network interfaces found",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			service := wifiPkg.New(tc.provider)
			result, err := service.GetNames()

			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorSubstring)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}