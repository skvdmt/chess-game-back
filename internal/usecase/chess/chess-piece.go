package chess

import (
	"github.com/google/uuid"
	"github.com/skvdmt/chess-game-back/internal/entities"
)

// ChessPiece Шахматная фигура.
type ChessPiece interface {
	// Идентификатор.
	Id() uuid.UUID
	// Имя.
	Name() string
	// На доске.
	OnBoard() bool
	// Позиция.
	Pos() *entities.Position
	// Позиция x.
	X() uint8
	// Позиция y.
	Y() uint8
	// Смена позиции.
	SetPos(pos *entities.Position, real bool)
	// Ходила.
	AlreadyMove() bool
	// Съесть.
	Eat()
	// Отмена съедения фигуры.
	UndoEat()
	// Rules Проверка хода на соответствие правилам хода фигуры..
	Rules(*entities.Position) bool
	// Направления.
	Directions() []entities.Direction
	// Максимальная дальность хода.
	MaxRemote() uint8
	// Возможные ходы.
	PossibleMoves() []*entities.Position
	// Установка возможных ходов.
	SetPossibleMoves([]*entities.Position)
}
