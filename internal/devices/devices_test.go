package devices

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

const layoutISO8601 = "2006-01-02 03:04:05Z07"

func TestChangeDeviceState(t *testing.T) {
	d := NewDevice("name01", "brand01")

	assert.Equal(t, d.State, Inactive)

	d.ChangeDeviceState(InUse)
	assert.Equal(t, d.State, InUse)

	d.ChangeDeviceState(Available)
	assert.Equal(t, d.State, InUse)
}

func TestChangeDeviceName_Success(t *testing.T) {
	d := NewDevice("name01", "brand01")

	newName := "name02"
	err := d.ChangeDeviceName(newName)

	if assert.NoError(t, err) {
		assert.Equal(t, d.Name, newName)
	}

}

func TestChangeDeviceName_Fail(t *testing.T) {
	d := NewDevice("name01", "brand01")
	d.State = InUse

	newName := "name02"
	err := d.ChangeDeviceName(newName)

	if assert.Error(t, err) {
		assert.Equal(t, "trying to change device %d name while in use", err)
	}
}

func TestChangeDeviceBrand_Success(t *testing.T) {
	d := NewDevice("name01", "brand01")

	newBrand := "brand02"
	err := d.ChangeDeviceBrand(newBrand)

	if assert.NoError(t, err) {
		assert.Equal(t, d.Brand, newBrand)
	}
}

func TestChangeDeviceBrand_Fail(t *testing.T) {
	d := NewDevice("name01", "brand01")
	d.State = InUse

	newBrand := "brand02"
	err := d.ChangeDeviceBrand(newBrand)

	if assert.Error(t, err) {
		assert.Equal(t, "trying to change device %d brand while in use", err)
	}
}

func TestIsDeviceInUse(t *testing.T) {
	d := NewDevice("name01", "brand01")

	b := d.IsDeviceInUse()
	assert.Equal(t, b, false)

	d.State = InUse
	b = d.IsDeviceInUse()
	assert.Equal(t, b, true)
}

func TestCreationTimeFormatted(t *testing.T) {
	d := NewDevice("name01", "brand01")

	d.CreatedAt = time.Date(2009, time.November, 10, 23, 1, 2, 0, time.FixedZone("UTC+2", 2))

	formattedTime := d.CreationTimeFormatted()

	expected := "2009-11-10T23:01:02+00:00"
	assert.Equal(t, expected, formattedTime)
}
