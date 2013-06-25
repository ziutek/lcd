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
	//buf [80 * 2 * 2]byte

	t time.Time
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

func (o *Bitbang) SetRW(b bool) error {
	if b {
		o.a |= o.rw
	} else {
		o.a &^= o.rw
	}
	o.buf[0] = o.a
	_, err := o.w.Write(o.buf[:1])
	return err
}

func (o *Bitbang) SetAUX(b bool) error {
	if b {
		o.a |= o.aux
	} else {
		o.a &^= o.aux
	}
	o.buf[0] = o.a
	_, err := o.w.Write(o.buf[:1])
	return err
}

func (o *Bitbang) wait() {
	if !o.t.IsZero() {
		d := o.t.Sub(time.Now())
		if d > 0 {
			time.Sleep(d)
		}
	}
}

func (o *Bitbang) setWait(d time.Duration) {
	o.t = time.Now().Add(d)
}

func (o *Bitbang) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}

	if len(data) == 1 {
		// One nibble: initialisation command
		o.wait()
		b := data[0] | o.a
		o.buf[0] = b
		o.buf[1] = b | o.e
		o.buf[2] = b
		_, err := o.w.Write(o.buf[:3])
		if err != nil {
			return 0, err
		}
		o.setWait(5 * time.Millisecond)
		return 1, nil
	}

	if len(data)%2 != 0 {
		panic("data length must be 1 or an even number")
	}

	for n := 0; n < len(data); n += 2 {
		// Multiple nibbles: regural commands
		o.wait()
		b0 := data[n]
		b1 := data[n+1]
		b := b0<<4 | b1&0x0f
		b0 |= o.a
		b1 |= o.a
		o.buf[0] = b0
		o.buf[1] = b0 | o.e
		o.buf[2] = b0
		o.buf[3] = b1
		o.buf[4] = b1 | o.e
		o.buf[5] = b1
		_, err := o.w.Write(o.buf[:])
		if err != nil {
			return n, err
		}
		if b < 4 {
			// "Clear display" or "Return home"
			o.setWait(16 * time.Millisecond)
		} else {
			// Other command
			o.setWait(40 * time.Microsecond)
		}
	}
	return len(data), nil
}

/*func (o *Bitbang) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}

	if len(data) == 1 {
		// One nibble: initialisation command
		b := data[0] | o.a
		o.buf[0] = b | o.e
		o.buf[1] = b
		_, err := o.w.Write(o.buf[:2])
		if err != nil {
			return 0, err
		}
		return 1, nil
	}

	if len(data)%2 != 0 {
		panic("data length must be 1 or an even number")
	}

	l := len(data)
	for len(data) > 0 {
		n := len(o.buf) / 2
		if n > len(data) {
			n = len(data)
		}
		k := 0
		for i := 0; i < n; i += 2 {
			b0 := data[i] | o.a
			b1 := data[i+1] | o.a
			o.buf[k] = b0 | o.e
			o.buf[k+1] = b0
			o.buf[k+2] = b1 | o.e
			o.buf[k+3] = b1
			k += 4
		}
		k, err := o.w.Write(o.buf[:k])
		if err != nil {
			return l - len(data) + k/2, err
		}
		data = data[n:]
	}
	return l, nil
}*/
