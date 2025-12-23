package wifi_test

import (
	"errors"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	mywifi "github.com/task-6/internal/wifi"
)

func TestWiFi(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockWiFi := new(WiFiHandle)
		service := mywifi.New(mockWiFi)
		mockWiFi.On("Interfaces").Return([]*wifi.Interface{{Name: "wlan0"}}, nil)

		names, err := service.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"wlan0"}, names)
	})

	t.Run("error", func(t *testing.T) {
		mockWiFi := new(WiFiHandle)
		service := mywifi.New(mockWiFi)
		mockWiFi.On("Interfaces").Return([]*wifi.Interface(nil), errors.New("err"))

		names, err := service.GetNames()
		assert.Error(t, err)
		assert.Nil(t, names)
	})
}
