package rest

import "devices_api/internal/devices"

type DevicesResponse struct {
	Body struct {
		Devices    []devices.Device `json:"devices"`
		StatusCode int              `json:"status_code"`
	} `json:"body"`
}

type SingleDeviceResponse struct {
	Body struct {
		Device     devices.Device `json:"device"`
		StatusCode int            `json:"status_code"`
	} `json:"body"`
}
