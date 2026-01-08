package wifi_test

import (
	wifi "github.com/mdlayher/wifi"
	mock "github.com/stretchr/testify/mock"
)

type WiFiHandle struct {
	mock.Mock
}

func (_m *WiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Interfaces")
	}

	var r0 []*wifi.Interface
	var r1 error

	if rf, ok := ret.Get(0).(func() ([]*wifi.Interface, error)); ok {
		return rf()
	}

	if rf, ok := ret.Get(0).(func() []*wifi.Interface); ok {
		r0 = rf()
	} else if ret.Get(0) != nil {
		r0, ok = ret.Get(0).([]*wifi.Interface)
		if !ok {
			panic("type assertion failed")
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func NewWiFiHandle(t interface {
	mock.TestingT
	Cleanup(func())
}) *WiFiHandle {
	mockObj := &WiFiHandle{}
	mockObj.Mock.Test(t)

	t.Cleanup(func() { mockObj.AssertExpectations(t) })

	return mockObj
}
