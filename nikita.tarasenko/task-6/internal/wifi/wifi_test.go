package wifi_test

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	wifiPkg "task-6/internal/wifi"
)

func mustMAC(s string) net.HardwareAddr {
	m, err := net.ParseMAC(s)
	if err != nil {
		panic(err)
	}
	return m
}

func TestDeviceCatalog_ListDeviceNames(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		cat := wifiPkg.NewDeviceCatalog(func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{
				{Name: "wlan0", HardwareAddr: mustMAC("00:11:22:33:44:55")},
				{Name: "wlx123", HardwareAddr: mustMAC("aa:bb:cc:dd:ee:ff")},
			}, nil
		})

		names, err := cat.ListDeviceNames()
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{"wlan0", "wlx123"}, names)
	})

	t.Run("empty source", func(t *testing.T) {
		cat := wifiPkg.NewDeviceCatalog(func() ([]*wifi.Interface, error) {
			return nil, nil
		})

		names, err := cat.ListDeviceNames()
		require.NoError(t, err)
		assert.Empty(t, names)
	})
}

func TestDeviceCatalog_ListDeviceNames_Error(t *testing.T) {
	t.Run("permission denied", func(t *testing.T) {
		cat := wifiPkg.NewDeviceCatalog(func() ([]*wifi.Interface, error) {
			return nil, &wifi.OpError{Err: wifi.ErrNotSupported}
		})

		_, err := cat.ListDeviceNames()
		assert.ErrorIs(t, err, wifiPkg.ErrWiFiAccessDenied)
	})

	t.Run("driver fault", func(t *testing.T) {
		cat := wifiPkg.NewDeviceCatalog(func() ([]*wifi.Interface, error) {
			return nil, &wifi.OpError{Err: errors.New("driver panic")}
		})

		_, err := cat.ListDeviceNames()
		assert.ErrorIs(t, err, wifiPkg.ErrWiFiDriverFault)
	})
}

func TestDeviceCatalog_DeviceMACIterator(t *testing.T) {
	cat := wifiPkg.NewDeviceCatalog(func() ([]*wifi.Interface, error) {
		return []*wifi.Interface{
			{Name: "wlan0", HardwareAddr: mustMAC("00:11:22:33:44:55")},
			{Name: "wlx123", HardwareAddr: mustMAC("aa:bb:cc:dd:ee:ff")},
		}, nil
	})

	var macs []net.HardwareAddr
	ch := cat.DeviceMACIterator()

	for mac := range ch {
		macs = append(macs, mac)
	}

	assert.Len(t, macs, 2)
	assert.Contains(t, macs, mustMAC("00:11:22:33:44:55"))
	assert.Contains(t, macs, mustMAC("aa:bb:cc:dd:ee:ff"))
}

func TestDeviceCatalog_DeviceMACIterator_Empty(t *testing.T) {
	cat := wifiPkg.NewDeviceCatalog(func() ([]*wifi.Interface, error) { return nil, nil })
	ch := cat.DeviceMACIterator()

	timeout := time.After(100 * time.Millisecond)
	select {
	case _, ok := <-ch:
		assert.False(t, ok, "channel should be closed")
	case <-timeout:
		t.Fatal("timeout waiting for channel close")
	}
}

func BenchmarkDeviceCatalog_ListDeviceNames(b *testing.B) {
	cat := wifiPkg.NewDeviceCatalog(func() ([]*wifi.Interface, error) {
		return []*wifi.Interface{
			{Name: "wlan0", HardwareAddr: mustMAC("00:11:22:33:44:55")},
			{Name: "wlx123", HardwareAddr: mustMAC("aa:bb:cc:dd:ee:ff")},
		}, nil
	})

	for range b.N {
		_, _ = cat.ListDeviceNames()
	}
}

func BenchmarkDeviceCatalog_Iterator(b *testing.B) {
	cat := wifiPkg.NewDeviceCatalog(func() ([]*wifi.Interface, error) {
		return []*wifi.Interface{
			{Name: "wlan0", HardwareAddr: mustMAC("00:11:22:33:44:55")},
			{Name: "wlx123", HardwareAddr: mustMAC("aa:bb:cc:dd:ee:ff")},
		}, nil
	})

	for range b.N {
		ch := cat.DeviceMACIterator()
		for range ch {
		}
	}
}
