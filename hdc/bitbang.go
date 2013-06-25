package hdc

import (
	"io"
	"time"
)

// Bitbang handles communication with HD44780 using some controller that
// implements clocked output of received data onto the D4-D7, E, RS pins (eg.
// FTDI USB chips in bitbang mode).
//
// It writes 3 bytes for one nibble to the provided io.Writer:
// - first:  with E bit unset, need >= 40 ns
// - second: with E bit set,   need >= 230 ns
// - thrid:  with E bit unset, need >= 10 ns
// Full E cycle need >= 500 ns.
// Baudrate that satisfy all this time constrains: 1 B / 230 ns = 4347826 B/s
//
// Additionaly Bitbang provides methods to controll AUX and R/W bits (if you
// want to use R/W bit as second AUX, the real R/W pin should be connected to
// the VSS).
type Bitbang struct {
	w          io.Writer
	e, rw, aux byte
	a          byte
	buf        [6]byte
}

// NewBitbang re
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
	if len(data) == 1 {
		// One nibble, initialisation
	}


	n := 0
	blen := len(o.buf) / 3
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
			o.buf[k+2] = b
			k += 3
		}
		k, err := o.w.Write(o.buf[:k])
		n += k / 4
		if err != nil {
			return n, err
		}
		data = data[l:]
	}
	time.Sleep(100 * time.Millisecond)
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
