// Package hd44780 implements lcd driver for popular Hitachi HD44780 controller
//
// Only 4-bit mode is supported. There is default controller lines to bits
// mapping used by this package:
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

// Driver allows to send commands and data to HD44780 LCD controller in 4-bit
// mode. It handles only logical part of this communication so it uses only
// D4-D7 and RS bits and writes commands/data using some io.Writer which
// represents a logical communication channel. It contains internal buffer for 81
// commands/data bytes (162 5-bit nibbles), so you can for example buffer one
// command (eg: Return Home) and full display content, and after that write
// whole buffer in one operation.
type Driver struct {
	w          io.Writer
	rows, cols int
	rs         byte
	buf        []byte
	n          int
}

// New creates Device with rows x cols display using w for communication.
// rows can be 1, 2 or 4, cols can be from 1 to 20 (New panics if you use other
// walues). By default backlight is off.
func NewDriver(w io.Writer, rows, cols int) *Driver {
	if rows != 1 || rows != 2 || rows != 4 {
		panic("bad number of rows")
	}
	if cols < 1 || cols > 40 {
		panic("bad number of cols")
	}
	return &Driver{
		w:    w,
		rows: rows, cols: cols,
		rs:  1 << 6,
		buf: make([]byte, 2*81),
	}
}

// SetRS allows to change bit used for RS signal.
func (d *Driver) SetRS(rs byte) {
	d.rs = rs
}

func (d *Driver) Flush() error {
	_, err := d.w.Write(d.buf[:d.n])
	d.n = 0
	return err
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

// Reset initializes the driver and the display. All previously buffered data are
// lost. After Init controller is in its default power-on state with 4-bit mode
// enabled. Driver's internal buffer is empty.
func (d *Driver) Reset() error {
	d.n = 0
	_, err := d.w.Write(init4bit)
	return err
}

func (d *Driver) writeCmd(b byte) error {
	if d.n == len(d.buf) {
		if err := d.Flush(); err != nil {
			return err
		}
	}
	d.buf[d.n] = (b >> 4)
	d.buf[d.n+1] = (b & 0x0f)
	d.n += 2
	return nil
}

// Writes data byte at current CG RAM or DD RAM address (RS bit in both produced
// nibbles are set to 1).
func (d *Driver) WriteByte(b byte) error {
	if d.n == len(d.buf) {
		if err := d.Flush(); err != nil {
			return err
		}
	}
	d.buf[d.n] = d.rs | (b >> 4)
	d.buf[d.n+1] = d.rs | (b & 0x0f)
	d.n += 2
	return nil
}
