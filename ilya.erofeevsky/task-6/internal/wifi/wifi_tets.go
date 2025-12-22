package wifi_test

import (
	"errors"
	"testing"

	mywifi "github.com/Ilya-Er0fick/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
)

func TestNetManager_GetActiveInterfaces(t *testing.T) {
	t.Parallel()

	t.Run("success_interfaces", func(t *testing.T) {
		mockWiFi := new(WiFiHandle)
		service := mywifi.New(mockWiFi)

		expected := []*wifi.Interface{
			{Name: "wlan_office"},
			{Name: "wlan_guest"},
		}

		mockWiFi.On("Interfaces").Return(expected, nil)

		names, err := service.GetActiveInterfaces()
		assert.NoError(t, err)
		assert.Len(t, names, 2)
		assert.Equal(t, "wlan_office", names[0])
		mockWiFi.AssertExpectations(t)
	})

	t.Run("hardware_failure", func(t *testing.T) {
		mockWiFi := new(WiFiHandle)
		service := mywifi.New(mockWiFi)

		mockWiFi.On("Interfaces").Return([]*wifi.Interface(nil), errors.New("hw error"))

		names, err := service.GetActiveInterfaces()
		assert.Error(t, err)
		assert.Nil(t, names)
	})
}
