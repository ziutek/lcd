package lcd

import (
	"image"
)

type Area interface {
	Show()
	Hide()
	Raise()
	Delete()
	// Move changes area position on the screen.
	Move(image.Point)
	// Resize changes area dimensions
	Resize(image.Point)
}

type ImageArea interface {
	Area
	Update(image.Image)
	SyncUpdate(image.Image)
}

type TextArea interface {
	Area
	Update([]rune)
	SyncUpdate([]rune)
}

type ImageScreen interface {
	// NewImageArea creates new image area on the screen.
	NewImageArea(image.Rectangle) ImageArea
}

type TextScreen interface {
	// NewTextArea creates new text area
	// pos contains image.Position on the screen (in screen cordinates), size
	// contains size of area in area coordinates (runes).
	NewTextArea(pos, size image.Point) TextArea
}
