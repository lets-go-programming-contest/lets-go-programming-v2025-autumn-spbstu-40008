package wifi_test

import (
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockWiFiHandle struct {
	mock.Mock
}

func (_m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	ret := _m.Called()

	var r0 []*wifi.Interface
	if val, ok := ret.Get(0).([]*wifi.Interface); ok {
		r0 = val
	} else if ret.Get(0) != nil {
		panic(fmt.Sprintf("unexpected type for return arg 0: %T", ret.Get(0)))
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		err := ret.Error(1)
		if err != nil {
			r1 = fmt.Errorf("mock error: %w", err)
		}
	}

	return r0, r1
}
