package lcd

type TextArea interface {
	Show()
	Hide()
	Raise()
	SetGeometry(x, y, width, height int)
	Update([]rune)
	SyncUpdate([]rune)
	Delete()
}

type TextScreen interface {
	// NewTextArea creates new text area of size (width, height) begins from
	// x column and y row. Created area is hidden and initialized using
	// transparent runes (rune(0) means transparent).
	NewTextArea(x, y, width, height int, sync bool) TextWindow
}
