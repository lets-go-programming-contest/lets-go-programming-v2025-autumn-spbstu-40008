package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/Ilya-Er0fick/task-6/internal/wifi"
	wifipkg "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockWiFi struct {
	mock.Mock
}

func (m *MockWiFi) Interfaces() ([]*wifipkg.Interface, error) {
	args := m.Called()
	return args.Get(0).([]*wifipkg.Interface), args.Error(1)
}

func TestWiFiService_Coverage(t *testing.T) {
	t.Parallel()
	
	errWifi := errors.New("wifi failure")
	mac, _ := net.ParseMAC("00:11:22:33:44:55")

	t.Run("GetAddresses_Success", func(t *testing.T) {
		m := new(MockWiFi)
		m.On("Interfaces").Return([]*wifipkg.Interface{{HardwareAddr: mac}}, nil)
		svc := wifi.New(m)
		res, err := svc.GetAddresses()
		require.NoError(t, err)
		require.Len(t, res, 1)
	})

	t.Run("GetAddresses_Error", func(t *testing.T) {
		m := new(MockWiFi)
		m.On("Interfaces").Return([]*wifipkg.Interface{}, errWifi)
		svc := wifi.New(m)
		_, err := svc.GetAddresses()
		require.Error(t, err)
	})

	t.Run("GetNames_Success", func(t *testing.T) {
		m := new(MockWiFi)
		m.On("Interfaces").Return([]*wifipkg.Interface{{Name: "wlan0"}}, nil)
		svc := wifi.New(m)
		res, err := svc.GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{"wlan0"}, res)
	})

	t.Run("GetNames_Error", func(t *testing.T) {
		m := new(MockWiFi)
		m.On("Interfaces").Return([]*wifipkg.Interface{}, errWifi)
		svc := wifi.New(m)
		_, err := svc.GetNames()
		require.Error(t, err)
	})
}
