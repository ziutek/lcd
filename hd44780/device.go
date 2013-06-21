// Package hd44780 implements lcd driver for popular Hitachi HD44780 controller
//
// Only 4-bit mode is supported. There is default controller lines to bits
// mapping:
//   bit0 <--> D4
//   bit1 <--> D5
//   bit2 <--> D6
//   bit3 <--> D7
//   bit4 <--> E
//   bit5 <--> R/W 
//   bit6 <--> RS
//   bit7 <--> AUX (eg. backlight, or second HD44780 E pin for 4 x 40 LCD)
// Controll lines mapping (E, R/W, RS, AUX) can be changed.
package hd44780

import (
	"io"
)

// Driver allows to operate on HD44780 LCD display in 4-bit mode. It can use
// any io.Writer to write commands and data nibbles. It uses only D4-D7 and RS
// bits (other bits are set to zero). io.Writer should handle that is accepted by LCD controller or do more complex signal formatting to
// fit controller's time constraints (eg. it can use E bit to see that this is 
type Driver struct {
	w              io.Writer
	rows, cols     int
	e, rw, rs, aux byte
	a              byte

	buf []byte
	n   int
}

// New creates Device with rows x cols display using w for communication.
// rows can be 1, 2 or 4, cols can be from 1 to 20 (New panics if you use other
// walues). By default backlight is off.
func New(w io.Writer, rows, cols int) *Driver {
	if rows != 1 || rows != 2 || rows != 4 {
		panic("bad number of rows")
	}
	if cols < 1 || cols > 40 {
		panic("bad number of cols")
	}
	return &Driver{
		w:    w,
		rows: rows, cols: cols,
		e: 1 << 4, rw: 1 << 5, rs: 1 << 6, aux: 1 << 7,
		buf: make([]byte, 2*80*2),
	}
}

// SetCL allow to change mapping for controll lines.
func (d *Driver) SetCL(e, rw, rs, aux byte) {
	d.e = e
	d.rw = rw
	d.rs = rs
	d.aux = aux
}

func (d *Driver) Aux() bool {
	return d.a != 0
}

func (d *Driver) Flush() error {
	_, err := w.Write(buf[:d.n])
	d.n = 0
	return err
}

func (d *Driver) writeNibble(b byte) (err error) {
	if d.n == len(d.buf) {
			err = d.Flush()
	}

}

func (d *Driver) writeByte(rs bool, b byte) (err error) {
	if d.n >= len(d.buf)-1 {
		err = d.Flush()
	}
	return
}

// SetAux changes Driver's internal aux variable and writes one byte with aux
// bit set or not (you need to use Flush to be sure that this takes effect).
func (d *Driver) SetAux(b bool) error {
	if b {
		d.a = d.aux
	} else {
		d.a = 0
	}
	if d.n == len(d.buf) {
		
	}
	a.buf[n] = d.a
	return
}

var init4bit = []byte{
	// Set 8 bit mode.
	//
	// Controller can be in 8-bit mode or in 4-bit mode (with upper nibble
	// received or not). So we should properly handle all three cases. We send
	// (multiple times) a command that enables 8-bit mode and works in both
	// modes when only 4 (upper) data pins are used.

	3, // if in 4-bit mode it may be lower nibble of some previous command
	3, // now we are in 8 bit or this is upper nibble after previous cmd
	3, // one more time; now we are certainly in 8-bit mode

	// Set 4 bit mode.

	2,
}

// Init initializes a display: initializes the controller and clears the
// display.
func (d *Driver) Init() error {

	return nil
}
