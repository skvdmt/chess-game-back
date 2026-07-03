package entities

// Позиция.
type Position struct {
	x uint8
	y uint8
}

// NewPosition Конструктор.
func NewPosition(x, y uint8) *Position {
	return &Position{
		x: x,
		y: y,
	}
}

// X коодината.
func (p *Position) X() uint8 {
	return p.x
}

// Y коодината.
func (p *Position) Y() uint8 {
	return p.y
}

// SetX Установка координаты x
func (p *Position) SetX(x uint8) {
	p.x = x
}

// SetY Установка координаты y
func (p *Position) SetY(y uint8) {
	p.y = y
}

// OnBoard Позиция на доске.
func (p *Position) OnBoard() bool {
	return p.x >= 1 && p.x <= 8 && p.y >= 1 && p.y <= 8
}
