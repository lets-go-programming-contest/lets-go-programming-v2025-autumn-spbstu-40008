package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	wifiPkg "github.com/Czeeen/lets-go-programming-v2025-autumn-spbstu-40008/nikita.prokopenko/task-6/internal/wifi"
)

type MockProvider struct {
	interfaces []*wifi.Interface
	err        error
}

func (m *MockProvider) Interfaces() ([]*wifi.Interface, error) {
	return m.interfaces, m.err
}

func createTestInterface(name, macAddress string) *wifi.Interface {
	mac, _ := net.ParseMAC(macAddress)
	return &wifi.Interface{
		Name:         name,
		HardwareAddr: mac,
	}
}

func TestNetworkService_GetAddresses(t *testing.T) {
	t.Parallel()

	t.Run("success with valid interfaces", func(t *testing.T) {
		provider := &MockProvider{
			interfaces: []*wifi.Interface{
				createTestInterface("eth0", "01:02:03:04:05:06"),
				createTestInterface("wlan0", "01:02:03:04:05:07"),
			},
		}
		service := wifiPkg.New(provider)
		result, err := service.GetAddresses()
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "01:02:03:04:05:06", result[0].String())
		assert.Equal(t, "01:02:03:04:05:07", result[1].String())
	})

	t.Run("error fetching interfaces", func(t *testing.T) {
		provider := &MockProvider{
			err: errors.New("interface error"),
		}
		service := wifiPkg.New(provider)
		result, err := service.GetAddresses()
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "failed to fetch interfaces")
	})

	t.Run("no interfaces", func(t *testing.T) {
		provider := &MockProvider{
			interfaces: []*wifi.Interface{},
		}
		service := wifiPkg.New(provider)
		result, err := service.GetAddresses()
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "no valid network interfaces found")
	})

	t.Run("interfaces without MAC addresses", func(t *testing.T) {
		provider := &MockProvider{
			interfaces: []*wifi.Interface{
				{Name: "lo", HardwareAddr: net.HardwareAddr{}},
			},
		}
		service := wifiPkg.New(provider)
		result, err := service.GetAddresses()
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "no valid network interfaces found")
	})

	t.Run("mix of valid and invalid interfaces", func(t *testing.T) {
		provider := &MockProvider{
			interfaces: []*wifi.Interface{
				createTestInterface("eth0", "01:02:03:04:05:06"),
				{Name: "lo", HardwareAddr: net.HardwareAddr{}},
				createTestInterface("wlan0", "01:02:03:04:05:07"),
			},
		}
		service := wifiPkg.New(provider)
		result, err := service.GetAddresses()
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "01:02:03:04:05:06", result[0].String())
		assert.Equal(t, "01:02:03:04:05:07", result[1].String())
	})
}

func TestNetworkService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success with valid interfaces", func(t *testing.T) {
		provider := &MockProvider{
			interfaces: []*wifi.Interface{
				{Name: "eth0"},
				{Name: "wlan0"},
			},
		}
		service := wifiPkg.New(provider)
		result, err := service.GetNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"eth0", "wlan0"}, result)
	})

	t.Run("error fetching interfaces", func(t *testing.T) {
		provider := &MockProvider{
			err: errors.New("permission denied"),
		}
		service := wifiPkg.New(provider)
		result, err := service.GetNames()
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "failed to fetch interfaces")
	})

	t.Run("no interfaces", func(t *testing.T) {
		provider := &MockProvider{
			interfaces: []*wifi.Interface{},
		}
		service := wifiPkg.New(provider)
		result, err := service.GetNames()
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "no valid network interfaces found")
	})

	t.Run("interfaces without names", func(t *testing.T) {
		provider := &MockProvider{
			interfaces: []*wifi.Interface{
				{Name: ""},
			},
		}
		service := wifiPkg.New(provider)
		result, err := service.GetNames()
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "no valid network interfaces found")
	})

	t.Run("mix of valid and invalid interfaces", func(t *testing.T) {
		provider := &MockProvider{
			interfaces: []*wifi.Interface{
				{Name: "eth0"},
				{Name: ""},
				{Name: "wlan0"},
			},
		}
		service := wifiPkg.New(provider)
		result, err := service.GetNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"eth0", "wlan0"}, result)
	})
}

func TestNetworkService_New(t *testing.T) {
	t.Parallel()

	provider := &MockProvider{}
	service := wifiPkg.New(provider)
	assert.NotNil(t, service)
}

func TestNetworkService_AllEmptyInterfaces(t *testing.T) {
	t.Parallel()

	provider := &MockProvider{
		interfaces: []*wifi.Interface{
			{Name: "", HardwareAddr: net.HardwareAddr{}},
			{Name: "", HardwareAddr: net.HardwareAddr{}},
		},
	}
	service := wifiPkg.New(provider)
	
	result, err := service.GetAddresses()
	assert.Nil(t, result)
	assert.ErrorContains(t, err, "no valid network interfaces found")
	
	result2, err2 := service.GetNames()
	assert.Nil(t, result2)
	assert.ErrorContains(t, err2, "no valid network interfaces found")
}

func TestNetworkService_SingleValidInterface(t *testing.T) {
	t.Parallel()

	provider := &MockProvider{
		interfaces: []*wifi.Interface{
			createTestInterface("eth0", "01:02:03:04:05:06"),
		},
	}
	service := wifiPkg.New(provider)
	
	result, err := service.GetAddresses()
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "01:02:03:04:05:06", result[0].String())
	
	result2, err2 := service.GetNames()
	require.NoError(t, err2)
	assert.Equal(t, []string{"eth0"}, result2)
}