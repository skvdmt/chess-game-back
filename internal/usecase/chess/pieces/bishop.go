package pieces

import (
	"github.com/google/uuid"
	"github.com/skvdmt/chess-game-back/internal/entities"
	"github.com/skvdmt/chess-game-back/internal/model"
)

// Bishop Офицер.
type Bishop struct {
	General
}

// NewBishop Конструктор офицера.
func NewBishop(pos *entities.Position) *Bishop {
	return &Bishop{
		General: General{
			id:       uuid.New(),
			name:     model.BISHOP,
			Position: pos,
			onBoard:  true,
		},
	}
}

// Rules Проверка хода на соответствие правилам хода фигуры.
func (b *Bishop) Rules(pos *entities.Position) bool {
	return (pos.X() < b.X() && pos.Y() < b.Y() && b.X()-pos.X() == b.Y()-pos.Y()) ||
		(pos.X() < b.X() && pos.Y() > b.Y() && b.X()-pos.X() == pos.Y()-b.Y()) ||
		(pos.X() > b.X() && pos.Y() < b.Y() && pos.X()-b.X() == b.Y()-pos.Y()) ||
		(pos.X() > b.X() && pos.Y() > b.Y() && pos.X()-b.X() == pos.Y()-b.Y())
}

// Направления.
func (b *Bishop) Directions() []entities.Direction {
	return []entities.Direction{
		entities.TopRight,
		entities.RightBottom,
		entities.BottomLeft,
		entities.LeftTop,
	}
}

// Максимальная дальность хода.
func (b *Bishop) MaxRemote() uint8 {
	return 8
}
