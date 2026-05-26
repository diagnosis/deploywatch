package handlers

import devicetokenstore "github.com/diagnosis/deploywatchv2/internal/store/devicetoken"

type DeviceTokenHandler struct {
	store devicetokenstore.DeviceTokenStore
}

func NewDeviceTokenHandler(store devicetokenstore.DeviceTokenStore) *DeviceTokenHandler {
	return &DeviceTokenHandler{store: store}
}
