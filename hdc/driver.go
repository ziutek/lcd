package hdc

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

// Push replaces current screen content.
// buf should contain new screen content line by line.
func (d *Driver) Refresh(buf []rune) error {
	if d.dev.rows == 1 {
	}
	return nil
}

func runesToBytes(bs []byte, rs []rune) {
	for i, r := range rs {
		if r >= 127 {
			bs[i] = '|'
		} else {
			bs[i] = byte(r)
		}
	}
}
