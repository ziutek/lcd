package hdc

import ()

// Driver implements lcd.Driver interface
type Driver struct {
	dev    *Device
	screen []byte
}

func NewDriver(dev *Device) *Driver {
	d := Driver{dev: dev}
	if d.dev.rows == 4 {
		d.screen = make([]byte, d.dev.rows*d.dev.cols)
	}
	return &d
}

// Size returns screen size
func (d *Driver) Size() (int, int) {
	return d.dev.rows, d.dev.cols
}

// Push replaces current scren content. screen should contain new screen
// content line by line.
func (d *Driver) Push(screen []byte) error {
	if d.dev.rows != 4 {
		_, err := d.dev.w.Write(screen)
		return err
	}
	w := d.dev.cols
	copy(d.screen, screen[:w])
	copy(d.screen[2*w:], screen[w:2*w])
	copy(d.screen[w:], screen[2*w:3*w])
	copy(d.screen[3*w:], screen[3*w:4*w])
	_, err := d.dev.w.Write(d.screen)
	return err
}
