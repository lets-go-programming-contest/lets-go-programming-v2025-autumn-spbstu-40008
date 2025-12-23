package wifi

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockWiFiClient struct {
	mock.Mock
}

func (m *mockWiFiClient) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*wifi.Interface), args.Error(1)
}

var errSim = errors.New("mocked interface error")

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	cases := []struct {
		title        string
		setupMock    func(*mockWiFiClient)
		wantAddrs    []net.HardwareAddr
		expectErr    bool
		errSubstring string
	}{
		{
			title: "returns hardware addresses",
			setupMock: func(m *mockWiFiClient) {
				hw, _ := net.ParseMAC("11:22:33:44:55:66")
				m.On("Interfaces").Return([]*wifi.Interface{{HardwareAddr: hw}}, nil)
			},
			wantAddrs: []net.HardwareAddr{
				{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
			},
			expectErr: false,
		},
		{
			title: "interfaces call fails",
			setupMock: func(m *mockWiFiClient) {
				m.On("Interfaces").Return(nil, errSim)
			},
			wantAddrs:    nil,
			expectErr:    true,
			errSubstring: "interface enumeration failed",
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()

			mockClient := new(mockWiFiClient)
			tc.setupMock(mockClient)

			service := New(mockClient) // ← правильно: New(w WiFi) WiFiService
			addrs, err := service.GetAddresses()

			if tc.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errSubstring)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantAddrs, addrs)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	cases := []struct {
		title        string
		setupMock    func(*mockWiFiClient)
		wantNames    []string
		expectErr    bool
		errSubstring string
	}{
		{
			title: "returns interface names",
			setupMock: func(m *mockWiFiClient) {
				m.On("Interfaces").Return([]*wifi.Interface{
					{Name: "wlp2s0"},
					{Name: "docker0"},
				}, nil)
			},
			wantNames: []string{"wlp2s0", "docker0"},
			expectErr: false,
		},
		{
			title: "interfaces call fails",
			setupMock: func(m *mockWiFiClient) {
				m.On("Interfaces").Return(nil, errSim)
			},
			wantNames:    nil,
			expectErr:    true,
			errSubstring: "failed to list interface names",
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()

			mockClient := new(mockWiFiClient)
			tc.setupMock(mockClient)

			service := New(mockClient)
			names, err := service.GetNames()

			if tc.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errSubstring)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantNames, names)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
