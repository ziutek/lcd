package lcd

import (
	"container/list"
)

type TextDriver interface {
	Size() (width, height int)
	Referesh([]rune) error
}

type TextDisplay struct {
	drv     TextDriver
	windows *list.List
}

func NewTextDisplay(drv TextDriver) *TextDisplay {
	return &TextDisplay{drv: drv}
}
