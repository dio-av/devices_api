package devices

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExist     = errors.New("row does not exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
	ErrDeviceInUse  = errors.New("cannot delete device while in use state")
)

// CreateDevice represents the model to create a new device.
type CreateDevice struct {
	Name  string      `json:"name"`
	Brand string      `json:"brand"`
	State DeviceState `json:"state"`
}

// Repository represents the interface contract for the Repository design pattern
// for the devices to implement any relational database.
type Repository interface {
	Writer
	Reader
	Service
}

// Writer represents the behaviour for writing data to repository.
type Writer interface {
	Create(ctx context.Context, cd CreateDevice) (*Device, error)
	Update(ctx context.Context, d Device) (sql.Result, error)
	Delete(ctx context.Context, d Device) (sql.Result, error)
}

// Reader represents the behaviour for reading data from repository.
type Reader interface {
	GetById(ctx context.Context, id int64) (*Device, error)
	GetByBrand(ctx context.Context, b string) ([]Device, error)
	GetByState(ctx context.Context, s DeviceState) ([]Device, error)
	All(ctx context.Context) ([]Device, error)
}

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}
