package rest

import "devices_api/internal/devices"

type DevicesResponse struct {
	Devices    []devices.Device `json:"devices"`
	StatusCode int              `json:"status_code"`
}

type SingleDeviceResponse struct {
	Device     devices.Device `json:"device"`
	StatusCode int            `json:"status_code"`
}
