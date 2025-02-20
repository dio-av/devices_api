package rest

import "devices_api/internal/devices"

// TODO:
// Request and Response payloads for the REST api.

// TODO:
// Request payload for Device data model.

// TODO:
// Response payload for the Device data model.
type DevicesResponse struct {
	device     []devices.Device
	statusCode int
}

type SingleDeviceResponse struct {
	device     devices.Device
	statusCode int
}
