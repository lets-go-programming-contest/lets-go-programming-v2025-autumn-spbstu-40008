package wifi_mocks

import (
    "github.com/stretchr/testify/mock"
    "github.com/mdlayher/wifi"
)

type WiFiHandle struct {
    mock.Mock
}

func (m *WiFiHandle) Interfaces() ([]*wifi.Interface, error) {
    args := m.Called()
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]*wifi.Interface), args.Error(1)
}

func (m *WiFiHandle) StationInfo(ifi *wifi.Interface) (*wifi.StationInfo, error) {
    args := m.Called(ifi)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*wifi.StationInfo), args.Error(1)
}
