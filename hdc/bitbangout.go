package hdc

import (
	"io"
	"log"
	"time"
)

// BitbangOut handles communication with HD44780 using some controller that
// implements clocked output of received data onto the D4-D7, E, RS pins (eg.
// FTDI USB chips in bitbang mode).
//
// Additionaly it provides methods to controll AUX and R/W bits (if you want to
// use R/W bit as second AUX, the real R/W pin should be connected to the VSS).
type BitbangOut struct {
	w          io.Writer
	e, rw, aux byte
	a          byte
	buf        []byte
}

func NewBitbangOut(w io.Writer) *BitbangOut {
	return &BitbangOut{
		w:   w,
		e:   1 << 4,
		rw:  1 << 5,
		aux: 1 << 7,
	}
}

func (o *BitbangOut) SetMapping(e, rw, aux byte) {
	o.e = e
	o.rw = rw
	o.aux = aux
}

func (o *BitbangOut) WriteC(data []byte) (int, error) {
	n := 0
	blen := len(o.buf) / 4
	for len(data) > 0 {
		l := len(data)
		if l > blen {
			l = blen
		}
		k := 0
		for _, b := range data[:l] {
			b |= o.a
			o.buf[k] = b
			o.buf[k+1] = b | o.e
			o.buf[k+2] = b | o.e
			o.buf[k+3] = b
			k += 4
		}
		k, err := o.w.Write(o.buf[:k])
		for _, b := range o.buf[:k] {
			log.Printf("%02x %08b", b, b)
		}
		n += k / 4
		if err != nil {
			return n, err
		}
		data = data[l:]
	}
	return n, nil
}

func (o *BitbangOut) Write(data []byte) (int, error) {
	buf := make([]byte, 2)
	for n, b := range data {
		b |= o.a
		buf[0] = b | o.e
		buf[1] = b
		if _, err := o.w.Write(buf[:1]); err != nil {
			return n - 1, err
		}
		time.Sleep(10 * time.Millisecond)
		if _, err := o.w.Write(buf[1:]); err != nil {
			return n - 1, err
		}
		time.Sleep(10 * time.Millisecond)
	}
	return len(data), nil
}
