package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"

	wifipkg "github.com/LAffey26/task-6/internal/wifi"
)

var errInterfaces = errors.New("interfaces error")

func parseMAC(value string) net.HardwareAddr {
	mac, _ := net.ParseMAC(value)

	return mac
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		handle := &mockWiFiHandle{}
		service := wifipkg.New(handle)

		handle.On("Interfaces").
			Return([]*wifi.Interface{
				{
					Name:         "wlan0",
					HardwareAddr: parseMAC("00:11:22:33:44:55"),
				},
			}, nil).
			Once()

		result, err := service.GetAddresses()

		require.NoError(t, err)
		require.Equal(
			t,
			[]net.HardwareAddr{parseMAC("00:11:22:33:44:55")},
			result,
		)

		handle.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		handle := &mockWiFiHandle{}
		service := wifipkg.New(handle)

		handle.On("Interfaces").
			Return(nil, errInterfaces).
			Once()

		result, err := service.GetAddresses()

		require.Error(t, err)
		require.Contains(t, err.Error(), "getting interfaces:")
		require.Nil(t, result)

		handle.AssertExpectations(t)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		handle := &mockWiFiHandle{}
		service := wifipkg.New(handle)

		handle.On("Interfaces").
			Return([]*wifi.Interface{
				{Name: "wlan0"},
				{Name: "eth0"},
			}, nil).
			Once()

		result, err := service.GetNames()

		require.NoError(t, err)
		require.Equal(t, []string{"wlan0", "eth0"}, result)

		handle.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		handle := &mockWiFiHandle{}
		service := wifipkg.New(handle)

		handle.On("Interfaces").
			Return(nil, errInterfaces).
			Once()

		result, err := service.GetNames()

		require.Error(t, err)
		require.Contains(t, err.Error(), "getting interfaces:")
		require.Nil(t, result)

		handle.AssertExpectations(t)
	})
}
