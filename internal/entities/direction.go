package entities

// Direction Направление.
type Direction uint8

const (
	// Вверх.
	Top = iota
	// Вверх направо.
	TopRight
	// Направо.
	Right
	// Направо вниз.
	RightBottom
	// Вниз.
	Bottom
	// Вниз налево.
	BottomLeft
	// Налево.
	Left
	// Налево вверх.
	LeftTop
)
