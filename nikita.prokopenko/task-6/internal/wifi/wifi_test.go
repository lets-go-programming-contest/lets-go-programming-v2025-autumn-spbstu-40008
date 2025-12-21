package netif_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Czeeen/lets-go-programming-v2025-autumn-spbstu-40008/nikita.prokopenko/task-6/internal/netif"
)

var (
	errInterfaceError = errors.New("interface access error")
	errPermission     = errors.New("permission denied for interface access")
)

func TestNetworkService_GetHardwareAddresses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupMock      func(*netif.MockInterfaceHandler)
		expectedMACs   []string
		expectError    bool
		errorSubstring string
	}{
		{
			name: "valid MAC addresses retrieval",
			setupMock: func(m *netif.MockInterfaceHandler) {
				interfaces := []*wifi.Interface{
					netif.CreateTestInterfaceData("eth0", "01:23:45:67:89:ab"),
					netif.CreateTestInterfaceData("wlan0", "cd:ef:01:23:45:67"),
				}
				m.On("FetchInterfaces").Return(interfaces, nil).Once()
			},
			expectedMACs: []string{"01:23:45:67:89:ab", "cd:ef:01:23:45:67"},
		},
		{
			name: "interface fetch failure",
			setupMock: func(m *netif.MockInterfaceHandler) {
				m.On("FetchInterfaces").Return(nil, errInterfaceError).Once()
			},
			expectError:    true,
			errorSubstring: "failed to fetch network interfaces",
		},
		{
			name: "no valid MAC addresses",
			setupMock: func(m *netif.MockInterfaceHandler) {
				interfaces := []*wifi.Interface{
					{Name: "lo", HardwareAddr: net.HardwareAddr{}},
					{Name: "dummy", HardwareAddr: net.HardwareAddr{0x00}},
				}
				m.On("FetchInterfaces").Return(interfaces, nil).Once()
			},
			expectError:    true,
			errorSubstring: "no valid MAC addresses detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			mockHandler := new(netif.MockInterfaceHandler)
			service := netif.NewNetworkService(mockHandler)
			tt.setupMock(mockHandler)

			macs, err := service.GetHardwareAddresses()

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorSubstring)
				assert.Nil(t, macs)
			} else {
				require.NoError(t, err)
				require.Len(t, macs, len(tt.expectedMACs))
				for i, expected := range tt.expectedMACs {
					assert.Equal(t, expected, macs[i].String())
				}
			}
			
			mockHandler.ValidateExpectations(t)
		})
	}
}

func TestNetworkService_GetInterfaceIdentifiers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupMock      func(*netif.MockInterfaceHandler)
		expectedNames  []string
		expectError    bool
		errorSubstring string
	}{
		{
			name: "successful interface names retrieval",
			setupMock: func(m *netif.MockInterfaceHandler) {
				interfaces := []*wifi.Interface{
					{Name: "eth0"},
					{Name: "wlan1"},
					{Name: "docker0"},
				}
				m.On("FetchInterfaces").Return(interfaces, nil).Once()
			},
			expectedNames: []string{"eth0", "wlan1", "docker0"},
		},
		{
			name: "permission denied error",
			setupMock: func(m *netif.MockInterfaceHandler) {
				m.On("FetchInterfaces").Return(nil, errPermission).Once()
			},
			expectError:    true,
			errorSubstring: "failed to fetch network interfaces",
		},
		{
			name: "all interface names empty",
			setupMock: func(m *netif.MockInterfaceHandler) {
				interfaces := []*wifi.Interface{
					{Name: ""},
					{Name: "  "},
				}
				m.On("FetchInterfaces").Return(interfaces, nil).Once()
			},
			expectError:    true,
			errorSubstring: "all interface names are empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			mockHandler := new(netif.MockInterfaceHandler)
			service := netif.NewNetworkService(mockHandler)
			tt.setupMock(mockHandler)

			names, err := service.GetInterfaceIdentifiers()

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorSubstring)
				assert.Nil(t, names)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedNames, names)
			}
			
			mockHandler.ValidateExpectations(t)
		})
	}
}

func (m *netif.MockInterfaceHandler) CreateTestInterfaceData(name, macAddress string) *wifi.Interface {
	return netif.CreateTestInterfaceData(name, macAddress)
}
