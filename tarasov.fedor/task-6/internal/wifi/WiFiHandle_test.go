package wifi_test

import (
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockWiFiHandle struct {
	mock.Mock
}

func (_m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	ret := _m.Called()

	var r0 []*wifi.Interface
	if rf, ok := ret.Get(0).(func() []*wifi.Interface); ok {
		r0 = rf()
	} else if ret.Get(0) != nil {
		if val, ok := ret.Get(0).([]*wifi.Interface); ok {
			r0 = val
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
