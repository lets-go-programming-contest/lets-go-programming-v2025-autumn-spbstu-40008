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
	errFetch    = errors.New("retrieve failed")
	errAccess   = errors.New("access denied")
	errExpected = "failed to fetch interfaces"
)

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	mockGen := func(macs []string) []*wifi.Interface {
		list := make([]*wifi.Interface, 0, len(macs))
		for i, m := range macs {
			hw, _ := net.ParseMAC(m)
			list = append(list, &wifi.Interface{
				Index:        i,
				Name:         "test_dev",
				HardwareAddr: hw,
			})
		}
		return list
	}

	tests := []struct {
		name      string
		mockSetup func(*MockWiFiHandle)
		want      []net.HardwareAddr
		wantErr   string
	}{
		{
			name: "Success",
			mockSetup: func(m *MockWiFiHandle) {
				data := mockGen([]string{"00:11:22:33:44:55"})
				m.On("Interfaces").Return(data, nil).Once()
			},
			want: []net.HardwareAddr{
				{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
			},
		},
		{
			name: "Error",
			mockSetup: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return(nil, errFetch).Once()
			},
			wantErr: errExpected,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := &MockWiFiHandle{}
			svc := New(m)

			tt.mockSetup(m)

			got, err := svc.GetAddresses()

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
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
		want      []string
		wantErr   string
	}{
		{
			name: "Success",
			mockSetup: func(m *MockWiFiHandle) {
				data := []*wifi.Interface{
					{Name: "wlan0"},
					{Name: "eth1"},
				}
				m.On("Interfaces").Return(data, nil).Once()
			},
			want: []string{"wlan0", "eth1"},
		},
		{
			name: "Error",
			mockSetup: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return(nil, errAccess).Once()
			},
			wantErr: errExpected,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := &MockWiFiHandle{}
			svc := New(m)

			tt.mockSetup(m)

			got, err := svc.GetNames()

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			m.AssertExpectations(t)
		})
	}
}