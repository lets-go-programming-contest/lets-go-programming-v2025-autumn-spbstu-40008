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
	errTypeAssertion  = errors.New("type assertion failed")
)

type MockInterfaceSource struct {
	mock.Mock
}

func (m *MockInterfaceSource) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	interfaces, ok := args.Get(0).([]*wifi.Interface)
	if !ok {
		return nil, errTypeAssertion
	}
	return interfaces, args.Error(1)
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
				}
				m.On("Interfaces").Return(interfaces, nil).Once()
			},
			expectedMACs: []string{"01:23:45:67:89:ab"},
		},
		{
			name: "fetch error",
			mockSetup: func(m *MockInterfaceSource) {
				m.On("Interfaces").Return(nil, errInterfaceError).Once()
			},
			expectError:    true,
			errorSubstring: "failed to fetch interfaces",
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
			} else {
				require.NoError(t, err)
				require.Len(t, macs, len(tc.expectedMACs))
				assert.Equal(t, tc.expectedMACs[0], macs[0].String())
			}
			mockSource.AssertExpectations(t)
		})
	}
}

func TestNetworkManager_GetInterfaceNames(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name          string
		mockSetup     func(*MockInterfaceSource)
		expectedNames []string
		expectError   bool
	}{
		{
			name: "valid interface names",
			mockSetup: func(m *MockInterfaceSource) {
				interfaces := []*wifi.Interface{{Name: "wlan0"}}
				m.On("Interfaces").Return(interfaces, nil).Once()
			},
			expectedNames: []string{"wlan0"},
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
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedNames, names)
			}
			mockSource.AssertExpectations(t)
		})
	}
}