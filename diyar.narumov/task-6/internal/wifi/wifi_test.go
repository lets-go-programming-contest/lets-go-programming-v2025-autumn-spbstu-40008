package wifi_test

import (
	"errors"
	"net"
	"testing"

	wifipkg "github.com/narumov-diyar/task-6/internal/wifi"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errDeviceAccessDenied = errors.New("device access denied")

func TestGetAddressesSuccess(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)

	mockInterfaces := []*wifi.Interface{
		{Name: "wlp2s0", HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}},
		{Name: "en0", HardwareAddr: net.HardwareAddr{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}},
	}
	mockWiFi.On("Interfaces").Return(mockInterfaces, nil)

	service := wifipkg.New(mockWiFi)

	addrs, err := service.GetAddresses()
	require.NoError(t, err)
	assert.Len(t, addrs, 2)
	assert.Equal(t, mockInterfaces[0].HardwareAddr, addrs[0])
	assert.Equal(t, mockInterfaces[1].HardwareAddr, addrs[1])

	mockWiFi.AssertExpectations(t)
}

func TestGetAddressesEmpty(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)

	mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, nil)

	service := wifipkg.New(mockWiFi)

	addrs, err := service.GetAddresses()
	require.NoError(t, err)
	assert.Empty(t, addrs)

	mockWiFi.AssertExpectations(t)
}

func TestGetAddressesError(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)

	mockWiFi.On("Interfaces").Return(nil, errDeviceAccessDenied)

	service := wifipkg.New(mockWiFi)

	_, err := service.GetAddresses()
	require.Error(t, err)
	require.EqualError(t, err, "getting interfaces: device access denied")

	mockWiFi.AssertExpectations(t)
}

func TestGetNamesSuccess(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)

	mockInterfaces := []*wifi.Interface{
		{Name: "wlan", HardwareAddr: net.HardwareAddr{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC}},
		{Name: "Wi-Fi", HardwareAddr: net.HardwareAddr{0xDE, 0xAD, 0xBE, 0xEF, 0xCA, 0xFE}},
	}
	mockWiFi.On("Interfaces").Return(mockInterfaces, nil)

	service := wifipkg.New(mockWiFi)

	names, err := service.GetNames()
	require.NoError(t, err)
	assert.Len(t, names, 2)
	assert.Equal(t, mockInterfaces[0].Name, names[0])
	assert.Equal(t, mockInterfaces[1].Name, names[1])

	mockWiFi.AssertExpectations(t)
}

func TestGetNamesEmpty(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)

	mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, nil)

	service := wifipkg.New(mockWiFi)

	names, err := service.GetNames()
	require.NoError(t, err)
	assert.Empty(t, names)

	mockWiFi.AssertExpectations(t)
}

func TestGetNamesError(t *testing.T) {
	t.Parallel()

	mockWiFi := NewWiFiHandle(t)

	mockWiFi.On("Interfaces").Return(nil, errDeviceAccessDenied)

	service := wifipkg.New(mockWiFi)

	_, err := service.GetNames()
	require.Error(t, err)
	require.EqualError(t, err, "getting interfaces: device access denied")

	mockWiFi.AssertExpectations(t)
}
