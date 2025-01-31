package devices

import "fmt"

func (d *Device) ChangeDeviceState(ds DeviceState) error {
	d.State = ds
	return nil
}

func (d *Device) ChangeDeviceName(n string) error {
	if d.State == InUse {
		return fmt.Errorf("trying to change device %d name while in use", d.Id)
	}

	d.Name = n
	return nil
}

func (d *Device) ChangeDeviceBrand(b string) error {
	if d.State == InUse {
		return fmt.Errorf("trying to change device %d brand while in use", d.Id)
	}

	d.Brand = b
	return nil
}

func (d *Device) CreationTimeFormatted() string {
	return d.CreatedAt.Format(dateTimeApiLayout)
}

func (d *Device) IsDeviceInUse() bool {
	return d.State == InUse
}
