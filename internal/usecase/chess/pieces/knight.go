package pieces

import (
	"github.com/google/uuid"
	"github.com/skvdmt/chess-game-back/internal/entities"
	"github.com/skvdmt/chess-game-back/internal/model"
)

// Knight Конь.
type Knight struct {
	General
}

// NewKnight Конструктор коня.
func NewKnight(pos *entities.Position) *Knight {
	return &Knight{
		General: General{
			id:       uuid.New(),
			name:     model.KNIGHT,
			Position: pos,
			onBoard:  true,
		},
	}
}

// Rules Проверка хода на соответствие правилам хода фигуры.
func (k *Knight) Rules(pos *entities.Position) bool {
	return (k.X()+1 == pos.X() && k.Y()+2 == pos.Y()) ||
		(k.X()+2 == pos.X() && k.Y()+1 == pos.Y()) ||
		(k.X()-1 == pos.X() && k.Y()-2 == pos.Y()) ||
		(k.X()-2 == pos.X() && k.Y()-1 == pos.Y()) ||
		(k.X()+1 == pos.X() && k.Y()-2 == pos.Y()) ||
		(k.X()+2 == pos.X() && k.Y()-1 == pos.Y()) ||
		(k.X()-1 == pos.X() && k.Y()+2 == pos.Y()) ||
		(k.X()-2 == pos.X() && k.Y()+1 == pos.Y())
}

// Направления.
func (k *Knight) Directions() []entities.Direction {
	return []entities.Direction{}
}

// Максимальная дальность хода.
func (k *Knight) MaxRemote() uint8 {
	return 1
}
