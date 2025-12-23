package wifi_test

import (
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockWiFi struct {
	mock.Mock
}

func NewWiFi(t mock.TestingT) *MockWiFi {
	m := &MockWiFi{}
	m.Test(t)
	return m
}

func (m *MockWiFi) Interfaces() ([]*wifi.Interface, error) {
	ret := m.Called()

	var r0 []*wifi.Interface
	if rf, ok := ret.Get(0).(func() []*wifi.Interface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*wifi.Interface)
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
