package pieces

import (
	"github.com/google/uuid"
	"github.com/skvdmt/chess-game-back/internal/entities"
)

// General Общая реализация шахматной фигуры.
type General struct {
	id            uuid.UUID
	name          string
	onBoard       bool
	alreadyMove   bool
	possibleMoves []*entities.Position
	*entities.Position
	// team        *Team
	// ChessPiece
}

// Id Идентификатор.
func (g *General) Id() uuid.UUID {
	return g.id
}

// Name Имя.
func (g *General) Name() string {
	return g.name
}

// OnBoadr На доске.
func (g *General) OnBoard() bool {
	return g.onBoard
}

// Pos Позиция.
func (g *General) Pos() *entities.Position {
	return g.Position
}

// X Позиция x.
func (g *General) X() uint8 {
	return g.Position.X()
}

// Y Позиция y.
func (g *General) Y() uint8 {
	return g.Position.Y()
}

// SetPos Смена позиции.
func (g *General) SetPos(pos *entities.Position, real bool) {
	g.Position = pos
	if real {
		g.alreadyMove = true
	}
}

// AlreadyMove Ходила.
func (g *General) AlreadyMove() bool {
	return g.alreadyMove
}

// Eat Съесть.
func (g *General) Eat() {
	g.onBoard = false
}

// UndoEat Отмена съедения фигуры.
func (g *General) UndoEat() {
	g.onBoard = true
}

// PossibleMoves Возможные ходы.
func (g *General) PossibleMoves() []*entities.Position {
	return g.possibleMoves
}

// SetPossibleMoves Установка возможных ходов.
func (g *General) SetPossibleMoves(pm []*entities.Position) {
	g.possibleMoves = pm
}
