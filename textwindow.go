package lcd

import (
	"unicode/utf8"
)

type TextWindow struct {
	width, height int
	runes         []rune
	cx, cy        int // cursor
}

func NewTextWindow(width, height int) TextWindow {
	w := &TextWindow{
		width:  width,
		height: height,
		runes:  make([]rune, width*height),
	}
	w.Clear()
	return w
}

func (w *TextWindow) Width() int {
	return w.width
}

func (w *TextWindow) Height() int {
	return w.height
}

func (w *TextWindow) Close() error {
	w.runes = nil
	return nil
}

func (w *TextWindow) Runes() []rune {
	return w.runes
}

func (w *TextWindow) CPos() (x, y int) {
	return w.cx, w.cy
}

func (w *TextWindow) SetCPos(x, y int) {
	w.cx = x
	w.cy = y
}

func (w *TextWindow) cshift(shift int) {
	w.cx += shift
	dy := w.cx / w.width
	w.cx -= dy * w.width
	w.cy += dy
}

func (w *TextWindow) Write(s []byte) (int, error) {
	addr := w.cy*w.width + w.cx
	for len(s) > 0 {
		r, l := utf8.DecodeRune(s)
		switch r {
		case '\n':
			addr = (addr/w.width + 1) * w.width
		case '\r':
			addr = (addr / w.width) * w.width
		default:
			w.runes[addr] = r
			addr = (addr + 1) % len(w.runes)
		}
		s = s[l:]
	}
	w.cx = 0
	w.cy = 0
	w.cshift(addr)
	return len(s), nil
}

func (w *TextWindow) Clear() {
	for n := range w.runes {
		w.runes[n] = ' '
	}
}
