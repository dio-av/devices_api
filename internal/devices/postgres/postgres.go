package postgres

import (
	"context"
	"database/sql"
	"devices_api/internal/devices"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *service
)

type service struct {
	db *sql.DB
}

func NewRepository() devices.Repository {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{
		db: db,
	}

	return dbInstance
}

const createDevice = `-- name: CreateDevice :one
INSERT INTO devices (
  d_name,
  d_brand,
  d_state,
  created_at
) VALUES (
  $1, $2, $3, NOW()
) RETURNING d_name, d_brand, d_state, created_at
`

func (s *service) Create(ctx context.Context, cd devices.CreateDevice) (*devices.Device, error) {
	row := s.db.QueryRowContext(ctx, createDevice, cd.Name, cd.Brand, cd.State)

	var d devices.Device
	err := row.Scan(
		&d.Id,
		&d.Name,
		&d.Brand,
		&d.State,
		&d.CreatedAt,
	)
	if err != nil {
		return &devices.Device{}, err
	}

	return &d, nil
}

const getDeviceById = `SELECT id, d_name, d_brand, d_state, created_at FROM devices
WHERE id = $1 LIMIT 1`

func (s *service) GetById(ctx context.Context, id int64) (*devices.Device, error) {
	var d devices.Device

	row := s.db.QueryRowContext(ctx, getDeviceById, id)

	err := row.Scan(
		&d.Id,
		&d.Name,
		&d.Brand,
		&d.State,
		&d.CreatedAt,
	)
	if err != nil {
		return &devices.Device{}, err
	}

	return &d, nil
}

const getDevicesByBrand = `SELECT id, d_name, d_brand, d_state, created_at FROM devices
WHERE d_brand = $1`

func (s *service) GetByBrand(ctx context.Context, brand string) ([]devices.Device, error) {
	var dd []devices.Device

	rows, err := s.db.QueryContext(ctx, getDevicesByBrand, brand)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d devices.Device

		err = rows.Scan(&d)
		if err != nil {
			return []devices.Device{}, err
		}

		err = rows.Err()
		if err != nil {
			return []devices.Device{}, err
		}
		dd = append(dd, d)
	}

	return dd, nil
}

const getDevicesByState = `SELECT id, d_name, d_brand, d_state, created_at FROM devices
WHERE d_state = $1`

func (s *service) GetByState(ctx context.Context, state devices.DeviceState) ([]devices.Device, error) {
	var dd []devices.Device

	rows, err := s.db.QueryContext(ctx, getDevicesByBrand, int(state))
	if err != nil {
		return []devices.Device{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var d devices.Device

		err = rows.Scan(&d)
		if err != nil {
			return []devices.Device{}, err
		}

		err = rows.Err()
		if err != nil {
			return []devices.Device{}, err
		}
		dd = append(dd, d)
	}

	return dd, nil
}

const getAllDevices = `SELECT * FROM devices`

func (s *service) All(ctx context.Context) ([]devices.Device, error) {
	var dd []devices.Device

	rows, err := s.db.QueryContext(ctx, getAllDevices)
	if err != nil {
		return []devices.Device{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var d devices.Device

		err = rows.Scan(&d)
		if err != nil {
			return []devices.Device{}, err
		}

		err = rows.Err()
		if err != nil {
			return []devices.Device{}, err
		}
		dd = append(dd, d)
	}

	return dd, nil
}

const updateDevice = `UPDATE devices SET d_name = $1, d_brand = $2, d_state = $3 WHERE id = $4;`

func (s *service) Update(ctx context.Context, d devices.Device) (sql.Result, error) {
	if d.DeviceInUse() {
		return nil, errors.New("cannot update device while in use state")
	}

	result, err := s.db.Exec(updateDevice, d.Name, d.Brand, d.State)
	if err != nil {
		return nil, err
	}

	return result, nil
}

const deleteDevice = `DELETE FROM devices where id = $1`

func (s *service) Delete(ctx context.Context, d devices.Device) (sql.Result, error) {
	// TODO: Check via database if the device is in use
	if d.DeviceInUse() {
		return nil, errors.New("cannot delete device while in use state")
	}

	result, err := s.db.Exec(deleteDevice, d.Id)
	if err != nil {
		return nil, devices.ErrDeleteFailed
	}

	return result, nil
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("%s", fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}
