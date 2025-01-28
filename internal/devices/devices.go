// Package Devices provides the implementation for models and behaviour for
// the device toghether with the definitions of the interfaces to interact
// with a database.
package devices

import (
	"time"
)

// swagger:model

// Device represents the model for the device.
type Device struct {
	Id        int64       `json:"id"`
	Name      string      `json:"name"`
	Brand     string      `json:"brand"`
	State     DeviceState `json:"state"`
	CreatedAt *time.Time  `json:"created_at"`
}

type DeviceState int

const (
	Available DeviceState = iota
	InUse
	Inactive
)

const dateTimeApiLayout = time.RFC3339

func NewDevice(id int64, name string, brand string) *Device {
	return &Device{
		Id:    id,
		Name:  name,
		Brand: brand,
		State: Inactive,
	}
}
