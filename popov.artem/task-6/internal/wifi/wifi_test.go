package wifi

import (
	"errors"
	wifiext "github.com/mdlayher/wifi"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockProvider struct {
	interfaces []*wifiext.Interface
	err        error
}

func (m *mockProvider) FetchInterfaces() ([]*wifiext.Interface, error) {
	return m.interfaces, m.err
}

var errMock = errors.New("mock error")

func TestRetrieveMACAddresses(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		provider      *mockProvider
		expectedAddrs []net.HardwareAddr
		expectErr     bool
	}{
		{
			name: "valid MAC addresses",
			provider: &mockProvider{
				interfaces: []*wifiext.Interface{
					{HardwareAddr: mustParseMAC("00:11:22:33:44:55")},
				},
			},
			expectedAddrs: []net.HardwareAddr{
				mustParseMAC("00:11:22:33:44:55"),
			},
			expectErr: false,
		},
		{
			name: "fetch error",
			provider: &mockProvider{
				err: errMock,
			},
			expectedAddrs: nil,
			expectErr:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			service := NewNetworkService(tc.provider)
			addrs, err := service.RetrieveMACAddresses()

			if tc.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "unable to retrieve interfaces")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedAddrs, addrs)
			}
		})
	}
}

func TestRetrieveInterfaceNames(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		provider      *mockProvider
		expectedNames []string
		expectErr     bool
	}{
		{
			name: "two names",
			provider: &mockProvider{
				interfaces: []*wifiext.Interface{
					{Name: "wlan0"},
					{Name: "wlan1"},
				},
			},
			expectedNames: []string{"wlan0", "wlan1"},
			expectErr:     false,
		},
		{
			name: "error",
			provider: &mockProvider{
				err: errMock,
			},
			expectedNames: nil,
			expectErr:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			service := NewNetworkService(tc.provider)
			names, err := service.RetrieveInterfaceNames()

			if tc.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "error fetching interface list")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedNames, names)
			}
		})
	}
}

func mustParseMAC(s string) net.HardwareAddr {
	addr, err := net.ParseMAC(s)
	if err != nil {
		panic(err)
	}
	return addr
}
