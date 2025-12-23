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
	errFetch = errors.New("fetch error")
	errAuth  = errors.New("auth error")
)

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	errExpected := "fetch interfaces"

	genIfaces := func(macs []string) []*wifi.Interface {
		res := make([]*wifi.Interface, 0, len(macs))

		for i, m := range macs {
			addr, _ := net.ParseMAC(m)
			res = append(res, &wifi.Interface{
				Index:        i,
				Name:         "test0",
				HardwareAddr: addr,
			})
		}

		return res
	}

	cases := []struct {
		name      string
		setup     func(*MockWiFiHandle)
		want      []net.HardwareAddr
		wantError string
	}{
		{
			name: "Success",
			setup: func(m *MockWiFiHandle) {
				m.On("Interfaces").
					Return(genIfaces([]string{"00:11:22:33:44:55"}), nil).
					Once()
			},
			want: []net.HardwareAddr{
				{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
			},
		},
		{
			name: "Error",
			setup: func(m *MockWiFiHandle) {
				m.On("Interfaces").
					Return(nil, errFetch).
					Once()
			},
			wantError: errExpected,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := &MockWiFiHandle{}
			svc := New(m)
			tc.setup(m)

			got, err := svc.GetAddresses()

			if tc.wantError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantError)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}

			m.AssertExpectations(t)
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	errExpected := "fetch interfaces"

	cases := []struct {
		name      string
		setup     func(*MockWiFiHandle)
		want      []string
		wantError string
	}{
		{
			name: "Success",
			setup: func(m *MockWiFiHandle) {
				data := []*wifi.Interface{
					{Name: "wlan0"},
					{Name: "eth0"},
				}
				m.On("Interfaces").Return(data, nil).Once()
			},
			want: []string{"wlan0", "eth0"},
		},
		{
			name: "Error",
			setup: func(m *MockWiFiHandle) {
				m.On("Interfaces").Return(nil, errAuth).Once()
			},
			wantError: errExpected,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := &MockWiFiHandle{}
			svc := New(m)
			tc.setup(m)

			got, err := svc.GetNames()

			if tc.wantError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantError)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}

			m.AssertExpectations(t)
		})
	}
}
