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

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifipkg.Interface, error) {
	args := m.Called()
	return args.Get(0).([]*wifipkg.Interface), args.Error(1)
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()
	mac, _ := net.ParseMAC("00:11:22:33:44:55")

	tests := []struct {
		name    string
		mockRet []*wifipkg.Interface
		mockErr error
		want    []net.HardwareAddr
		wantErr bool
	}{
		{
			name:    "success",
			mockRet: []*wifipkg.Interface{{HardwareAddr: mac}},
			want:    []net.HardwareAddr{mac},
		},
		{
			name:    "error",
			mockRet: []*wifipkg.Interface{},
			mockErr: errors.New("wifi error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockWiFiHandle{}
			m.On("Interfaces").Return(tt.mockRet, tt.mockErr)

			svc := wifi.New(m)
			got, err := svc.GetAddresses()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
