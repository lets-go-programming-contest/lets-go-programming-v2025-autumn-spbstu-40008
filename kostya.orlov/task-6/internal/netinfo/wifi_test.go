package netinfo_test  // <-- ВАЖНО: добавить _test

import (
    "errors"
    "net"
    "testing"

    "github.com/mdlayher/wifi"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/task-6/internal/netinfo"  // <-- Теперь это нормальный импорт
)

// MockScanner остается здесь
type MockScanner struct {
    MockFn func() ([]*wifi.Interface, error)
}

func (m *MockScanner) Interfaces() ([]*wifi.Interface, error) {
    return m.MockFn()
}

func TestWiFiManager_FetchMACAddresses(t *testing.T) {
    tests := []struct {
        name        string
        mockFn      func() ([]*wifi.Interface, error)
        expected    []string
        expectError bool
        errorMsg    string
    }{
        {
            name: "success with multiple interfaces",
            mockFn: func() ([]*wifi.Interface, error) {
                return []*wifi.Interface{
                    {
                        Name:         "wlan0",
                        HardwareAddr: net.HardwareAddr{0x00, 0x1A, 0x2B, 0x3C, 0x4D, 0x5E},
                    },
                    {
                        Name:         "eth0",
                        HardwareAddr: net.HardwareAddr{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF},
                    },
                }, nil
            },
            expected:    []string{"00:1a:2b:3c:4d:5e", "aa:bb:cc:dd:ee:ff"},
            expectError: false,
        },
        {
            name: "success with empty interfaces",
            mockFn: func() ([]*wifi.Interface, error) {
                return []*wifi.Interface{}, nil
            },
            expected:    []string{},
            expectError: false,
        },
        {
            name: "scanner returns error",
            mockFn: func() ([]*wifi.Interface, error) {
                return nil, errors.New("scanner failed")
            },
            expected:    nil,
            expectError: true,
            errorMsg:    "network error",
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            mock := &MockScanner{MockFn: tc.mockFn}
            mgr := netinfo.NewWiFiManager(mock)  // <-- Используем полный путь
            
            res, err := mgr.FetchMACAddresses()
            
            if tc.expectError {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tc.errorMsg)
                assert.Nil(t, res)
            } else {
                require.NoError(t, err)
                assert.Equal(t, tc.expected, res)
            }
        })
    }
}

func TestNewWiFiManager(t *testing.T) {
    mock := &MockScanner{}
    mgr := netinfo.NewWiFiManager(mock)  // <-- Используем полный путь
    
    assert.NotNil(t, mgr)
    // Нельзя проверить приватное поле scanner напрямую
    // Но можно проверить через поведение
    res, err := mgr.FetchMACAddresses()
    assert.Error(t, err)  // mock вернет nil, nil по умолчанию
    assert.Nil(t, res)
}