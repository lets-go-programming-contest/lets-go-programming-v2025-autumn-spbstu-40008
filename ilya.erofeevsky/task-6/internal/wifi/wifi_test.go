package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	wifipkg "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"ilya.erofeevsky/task-6/internal/wifi"
)

var errWifiStatic = errors.New("wifi failure")

type MockWiFi struct {
	mock.Mock
}

func (m *MockWiFi) Interfaces() ([]*wifipkg.Interface, error) {
	args := m.Called()
	res, _ := args.Get(0).([]*wifipkg.Interface)
	err := args.Error(1)
	if err != nil {
		return res, fmt.Errorf("mock error: %w", err)
	}

	return res, nil
}

func TestWiFiService_Coverage(t *testing.T) {
	t.Parallel()

	mac, _ := net.ParseMAC("00:11:22:33:44:55")

	t.Run("GetAddresses_Success", func(t *testing.T) {
		t.Parallel()

		m := new(MockWiFi)

		m.On("Interfaces").Return([]*wifipkg.Interface{{HardwareAddr: mac}}, nil)

		svc := wifi.New(m)
		res, err := svc.GetAddresses()

		require.NoError(t, err)
		require.Len(t, res, 1)
	})

	t.Run("GetAddresses_Error", func(t *testing.T) {
		t.Parallel()

		m := new(MockWiFi)

		m.On("Interfaces").Return([]*wifipkg.Interface{}, errWifiStatic)

		svc := wifi.New(m)
		_, err := svc.GetAddresses()

		require.Error(t, err)
	})

	t.Run("GetNames_Success", func(t *testing.T) {
		t.Parallel()

		m := new(MockWiFi)

		m.On("Interfaces").Return([]*wifipkg.Interface{{Name: "wlan0"}}, nil)

		svc := wifi.New(m)
		res, err := svc.GetNames()

		require.NoError(t, err)
		require.Equal(t, []string{"wlan0"}, res)
	})

	t.Run("GetNames_Error", func(t *testing.T) {
		t.Parallel()

		m := new(MockWiFi)

		m.On("Interfaces").Return([]*wifipkg.Interface{}, errWifiStatic)

		svc := wifi.New(m)
		_, err := svc.GetNames()

		require.Error(t, err)
	})
}
