package wifi

import (
	"github.com/mdlayher/wifi"
)

type MockProvider struct {
	FetchResponse []*wifi.Interface
	FetchError    error
}

func (m *MockProvider) FetchInterfaces() ([]*wifi.Interface, error) {
	return m.FetchResponse, m.FetchError
}
