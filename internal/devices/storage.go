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
)

type CreateDevice struct {
	Name  string      `json:"name"`
	Brand string      `json:"brand"`
	State DeviceState `json:"state"`
}

type Repository interface {
	Writer
	Reader
	Service
}

type Writer interface {
	Create(ctx context.Context, cd CreateDevice) (*Device, error)
	Update(ctx context.Context, d Device) (sql.Result, error)
	Delete(ctx context.Context, d Device) (sql.Result, error)
}

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
