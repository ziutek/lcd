package hdc

import (
	"io"
	"time"
)

// Bitbang handles communication with HD44780 using some controller that
// implements clocked output of received data onto the D4-D7, E, RS pins (eg.
// FTDI USB chips in bitbang mode).
//
// Additionaly it provides methods to controll AUX and R/W bits (if you want to
// use R/W bit as second AUX, the real R/W pin should be connected to the VSS).
type Bitbang struct {
	w          io.Writer
	e, rw, aux byte
	a          byte
	buf        [80 * 2 * 2]byte
}

func NewBitbang(w io.Writer) *Bitbang {
	return &Bitbang{
		w:   w,
		e:   1 << 4,
		rw:  1 << 5,
		aux: 1 << 7,
	}
}

func (o *Bitbang) SetMapping(e, rw, aux byte) {
	o.e = e
	o.rw = rw
	o.aux = aux
}

func (o *Bitbang) Write(data []byte) (int, error) {
	n := 0
	blen := len(o.buf) / 2
	for len(data) != 0 {
		l := len(data)
		if l > blen {
			l = blen
		}
		k := 0
		for _, b := range data[:l] {
			b |= o.a
			o.buf[k] = b
			o.buf[k+1] = b | o.e
			k += 2
		}
		k, err := o.w.Write(o.buf[:k])
		n += k / 2
		if err != nil {
			return n, err
		}
		data = data[l:]
	}
	time.Sleep(5 * time.Millisecond)
	return n, nil
}

/*
func (o *Bitbang) WriteC(data []byte) (int, error) {
	buf := make([]byte, 4)
	for n := 0; n < len(data); n += 2 {
		b := data[n]
		buf[0] = b | o.e
		buf[1] = b
		b = data[n+1]
		buf[2] = b | o.e
		buf[3] = b
		if _, err := o.w.Write(buf); err != nil {
			return n - 1, err
		}
		time.Sleep(time.Millisecond)
	}
	return len(data), nil
}*/
