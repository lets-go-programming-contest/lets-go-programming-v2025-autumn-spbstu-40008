package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	wifiPkg "task-6/internal/wifi"
)

var (
	errDrv  = errors.New("ошибка драйвера")
	errPerm = errors.New("отказано в доступе")
	errType = errors.New("ошибка приведения типа")
)

// MockWiFiHandle имитирует WiFiHandle
type MockWiFiHandle struct {
	mock.Mock
}

// Interfaces реализует интерфейс WiFiHandle
func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	if args.Get(0) == nil {
		e := args.Error(1)
		if e != nil {
			return nil, fmt.Errorf("ошибка мока: %w", e)
		}

		return nil, nil
	}

	ifaces, ok := args.Get(0).([]*wifi.Interface)
	if !ok {
		e := args.Error(1)
		if e != nil {
			return nil, fmt.Errorf("ошибка приведения типа: %w", e)
		}

		return nil, errType
	}

	e := args.Error(1)
	if e != nil {
		return ifaces, fmt.Errorf("ошибка результата мока: %w", e)
	}

	return ifaces, nil
}

func createMockIfaces(macs []string) []*wifi.Interface {
	res := make([]*wifi.Interface, 0, len(macs))

	for idx, mac := range macs {
		hw, parseErr := net.ParseMAC(mac)
		if parseErr != nil {
			continue
		}

		iface := &wifi.Interface{
			Index:        idx,
			Name:         "wlan" + string(rune('0'+idx)),
			HardwareAddr: hw,
		}
		res = append(res, iface)
	}

	return res
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name      string
		macs      []string
		mockErr   error
		expectErr bool
	}{
		{
			name:      "успешное получение MAC-адресов",
			macs:      []string{"00:11:22:33:44:55", "AA:BB:CC:DD:EE:FF"},
			mockErr:   nil,
			expectErr: false,
		},
		{
			name:      "список интерфейсов пуст",
			macs:      []string{},
			mockErr:   nil,
			expectErr: false,
		},
		{
			name:      "ошибка интерфейсов",
			macs:      nil,
			mockErr:   errDrv,
			expectErr: true,
		},
		{
			name:      "некорректный MAC-адрес",
			macs:      []string{"00:11:22:33:44:55", "неправильный-mac"},
			mockErr:   nil,
			expectErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHandle := &MockWiFiHandle{}

			var ifaces []*wifi.Interface
			if tc.macs != nil {
				ifaces = createMockIfaces(tc.macs)
			}

			mockHandle.On("Interfaces").Return(ifaces, tc.mockErr)

			service := wifiPkg.New(mockHandle)

			addrs, err := service.GetAddresses()

			if tc.expectErr {
				require.Error(t, err)
				assert.Nil(t, addrs)
			} else {
				require.NoError(t, err)

				validCount := 0

				for _, m := range tc.macs {
					if _, e := net.ParseMAC(m); e == nil {
						validCount++
					}
				}

				if validCount == 0 {
					assert.Empty(t, addrs)
				} else {
					assert.Len(t, addrs, validCount)

					for _, a := range addrs {
						assert.NotNil(t, a)
					}
				}
			}

			mockHandle.AssertExpectations(t)
		})
	}
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name      string
		names     []string
		mockErr   error
		expectErr bool
	}{
		{
			name:      "успешное получение имен",
			names:     []string{"wlan0", "wlan1", "eth0"},
			mockErr:   nil,
			expectErr: false,
		},
		{
			name:      "ошибка интерфейсов",
			names:     nil,
			mockErr:   errPerm,
			expectErr: true,
		},
		{
			name:      "нет интерфейсов",
			names:     []string{},
			mockErr:   nil,
			expectErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHandle := &MockWiFiHandle{}

			var ifaces []*wifi.Interface

			for i, n := range tc.names {
				iface := &wifi.Interface{
					Index: i,
					Name:  n,
				}
				ifaces = append(ifaces, iface)
			}

			mockHandle.On("Interfaces").Return(ifaces, tc.mockErr)

			service := wifiPkg.New(mockHandle)

			names, err := service.GetNames()

			if tc.expectErr {
				require.Error(t, err)
				assert.Nil(t, names)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.names, names)
			}

			mockHandle.AssertExpectations(t)
		})
	}
}

func BenchmarkGetAddresses(b *testing.B) {
	mockHandle := &MockWiFiHandle{}
	macs := []string{"00:11:22:33:44:55", "AA:BB:CC:DD:EE:FF"}
	ifaces := createMockIfaces(macs)

	mockHandle.On("Interfaces").Return(ifaces, nil).Times(b.N)

	s := wifiPkg.New(mockHandle)

	b.ResetTimer()

	for range b.N {
		_, _ = s.GetAddresses()
	}
}

func BenchmarkGetNames(b *testing.B) {
	mockHandle := &MockWiFiHandle{}
	names := []string{"wlan0", "wlan1", "eth0"}

	ifaces := make([]*wifi.Interface, 0, len(names))

	for idx, n := range names {
		iface := &wifi.Interface{
			Index: idx,
			Name:  n,
		}
		ifaces = append(ifaces, iface)
	}

	mockHandle.On("Interfaces").Return(ifaces, nil).Times(b.N)

	s := wifiPkg.New(mockHandle)

	b.ResetTimer()

	for range b.N {
		_, _ = s.GetNames()
	}
}
