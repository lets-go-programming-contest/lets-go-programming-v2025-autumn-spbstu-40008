package wifi

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errFetch  = errors.New("failed to fetch")
	errAccess = errors.New("access denied")
)

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	// Хелпер для создания тестовых данных
	makeIfaces := func(macs []string) []*wifi.Interface {
		list := make([]*wifi.Interface, 0, len(macs))
		for i, m := range macs {
			hw, _ := net.ParseMAC(m)
			list = append(list, &wifi.Interface{
				Index:        i,
				Name:         "wlan_test",
				HardwareAddr: hw,
			})
		}
		return list
	}

	tests := []struct {
		name      string
		mockSetup func(*MockWiFiHandle)
		expected  []net.HardwareAddr
		wantErr   string
	}{
		{
			name: "Success case",
			mockSetup: func(m *MockWiFiHandle) {
				data := makeIfaces([]string{"00:11:22:33:44:55"})
				m.On("Interfaces").Return(data, nil).Once()
			},
			expected: []net.HardwareAddr{
				{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
			},
		},
		{
			name: "Error case",
			mockSetup: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return(nil, errFetch).Once()
			},
			wantErr: "interface retrieval failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := &MockWiFiHandle{}
			svc := New(m)

			tc.mockSetup(m)

			result, err := svc.GetAddresses()

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}

			m.AssertExpectations(t)
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		mockSetup func(*MockWiFiHandle)
		expected  []string
		wantErr   string
	}{
		{
			name: "Success case",
			mockSetup: func(m *MockWiFiHandle) {
				data := []*wifi.Interface{
					{Name: "wlan0"},
					{Name: "eth0"},
				}
				m.On("Interfaces").Return(data, nil).Once()
			},
			expected: []string{"wlan0", "eth0"},
		},
		{
			name: "Error case",
			mockSetup: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return(nil, errAccess).Once()
			},
			wantErr: "interface retrieval failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := &MockWiFiHandle{}
			svc := New(m)

			tc.mockSetup(m)

			result, err := svc.GetNames()

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}

			m.AssertExpectations(t)
		})
	}
}