package lcd

import (
	"io"
	"unicode/utf8"
)

type Window interface {
	io.WriteCloser

	Width() int
	Height() int
}

type TextWindow interface {
	Window

	Runes() []rune
	CPos() (x, y int)
	SetCPos(x, y int)
	Clear()
}

type textWindow struct {
	// Content
	width, height int
	runes         []rune
	// Cursor
	cx, cy int
}

func NewTextWindow(width, height int) TextWindow {
	w := &textWindow{
		width:  width,
		height: height,
		runes:  make([]rune, width*height),
	}
	w.Clear()
	return w
}

func (w *textWindow) Width() int {
	return w.width
}

func (w *textWindow) Height() int {
	return w.height
}

func (w *textWindow) Close() error {
	w.runes = nil
	return nil
}

func (w *textWindow) Runes() []rune {
	return w.runes
}

func (w *textWindow) CPos() (x, y int) {
	return w.cx, w.cy
}

func (w *textWindow) SetCPos(x, y int) {
	w.cx = x
	w.cy = y
}

func (w *textWindow) cshift(shift int) {
	w.cx += shift
	dy := w.cx / w.width
	w.cx -= dy * w.width
	w.cy += dy
}

func (w *textWindow) Write(s []byte) (int, error) {
	addr := w.cy*w.width + w.cx
	for len(s) > 0 {
		r, l := utf8.DecodeRune(s)
		switch r {
		case '\n':
			addr = (addr/w.width + 1) * w.width
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

func (w *textWindow) Clear() {
	for n := range w.runes {
		w.runes[n] = ' '
	}
}
