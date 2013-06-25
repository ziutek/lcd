// Package hdc implements lcd.Device for popular Hitachi HD44780 controller
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
package hdc

import (
	"io"
)

// Device allows to send commands and data to HD44780 LCD controller in 4-bit
// mode. It handles only logical part of this communication so it uses only
// D4-D7 and RS bits and writes commands/data using some io.Writer which
// represents a logical communication channel.
//
// All commands are written to the provided io.Writer as two bytes (multiple
// commands can be written using one Write call) with exception of
// initialisation commands that are always written as one byte (only one
// initialisation command is written at a time).
type Device struct {
	w          io.Writer
	rows, cols int
	rs         byte
	buf        [80 * 2]byte
}

// NewDevice creates Device with rows x cols display using w for communication.
// rows can be 1, 2 or 4, cols can be from 1 to 20 (NewDevice panics if you use
// other walues).
func NewDevice(w io.Writer, rows, cols int) *Device {
	if rows != 1 && rows != 2 && rows != 4 {
		panic("bad number of rows")
	}
	if cols < 1 || cols > 40 {
		panic("bad number of cols")
	}
	return &Device{
		w:    w,
		rows: rows, cols: cols,
		rs: 1 << 6,
	}
}

// SetMapping allows to change bit used for RS signal.
func (d *Device) SetMapping(rs byte) {
	d.rs = rs
}

func (d *Device) writeNibble(b byte) error {
	d.buf[0] = b
	_, err := d.w.Write(d.buf[:1])
	return err
}

func (d *Device) writeCmd(b byte) error {
	d.buf[0] = b >> 4
	d.buf[1] = b & 0x0f
	_, err := d.w.Write(d.buf[:2])
	return err
}

func (d *Device) ClearDisplay() error {
	return d.writeCmd(0x01)
}

func (d *Device) ReturnHome() error {
	return d.writeCmd(0x02)
}

type EntryMode byte

const (
	DecrMode  EntryMode = 0
	IncrMode  EntryMode = 1 << 1
	ShiftMode EntryMode = 1
)

func (d *Device) SetEntryMode(f EntryMode) error {
	return d.writeCmd(byte(0x04 | f&0x03))
}

type Display byte

const (
	DisplayOff Display = 0
	DisplayOn  Display = 1 << 2
	CursorOff  Display = 0
	CursorOn   Display = 1 << 1
	BlinkOff   Display = 0
	BlinkOn    Display = 1
)

func (d *Device) SetDisplay(f Display) error {
	return d.writeCmd(byte(0x08 | f&7))
}

type Shift byte

const (
	ShiftCuror  Shift = 0
	ShiftScreen Shift = 1 << 3
	ShiftLeft   Shift = 0
	ShiftRight  Shift = 1 << 2
)

func (d *Device) SetShift(f Shift) error {
	return d.writeCmd(byte(0x10 | f&0xc))
}

type Function byte

const (
	OneLine  Function = 0
	TwoLines Function = 1 << 3
	Font5x8  Function = 0
	Font5x10 Function = 1 << 2
)

func (d *Device) SetFunction(f Function) error {
	return d.writeCmd(byte(0x20 | f&0x0f))
}

func (d *Device) SetCGRAMAddr(addr int) error {
	return d.writeCmd(0x40 | byte(addr)&0x3f)
}

func (d *Device) SetDDRAMAddr(addr int) error {
	return d.writeCmd(0x80 | byte(addr)&0x7f)
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

// Init initializes the driver and the display. All previously buffered data
// are lost. After Init controller should be in the following state:
// - 4-bit mode,
// - one line display if rows == 1, two line display otherwise,
// - 5x8 font,
// - display off, cursor off, blink off,
// - increment mode,
// - display cleared.
func (d *Device) Init() error {
	var err error
	for _, b := range init4bit {
		if err = d.writeNibble(b); err != nil {
			return err
		}
	}
	// Some controller models may require to use SetFunction before any other
	// instuction.
	f := Font5x8
	if d.rows != 1 {
		f |= TwoLines
	}
	err = d.SetFunction(f)
	if err != nil {
		return err
	}
	err = d.SetDisplay(DisplayOff | CursorOff | BlinkOff)
	if err != nil {
		return err
	}
	err = d.SetEntryMode(IncrMode)
	if err != nil {
		return err
	}
	return d.ClearDisplay()
}

// Write writes buf starting from current CG RAM or DD RAM address.
func (d *Device) Write(data []byte) (int, error) {
	for _, b := range data {
		d.WriteByte(b)
	}
	return len(data), nil

	n := 0
	blen := len(d.buf) / 2
	for len(data) != 0 {
		l := len(data)
		if l > blen {
			l = blen
		}
		k := 0
		for _, b := range data[:l] {
			d.buf[k] = d.rs | b>>4
			d.buf[k+1] = d.rs | b&0x0f
			k += 2
		}
		k, err := d.w.Write(d.buf[:k])
		n += k / 2
		if err != nil {
			return n, err
		}
		data = data[l:]
	}
	return n, nil
}

// Writes byte at current CG RAM or DD RAM address.
func (d *Device) WriteByte(b byte) error {
	d.buf[0] = d.rs | b>>4
	d.buf[1] = d.rs | b&0x0f
	_, err := d.w.Write(d.buf[:2])
	return err
}
