package pieces

import (
	"github.com/google/uuid"
	"github.com/skvdmt/chess-game-back/internal/entities"
	"github.com/skvdmt/chess-game-back/internal/model"
)

// King Король.
type King struct {
	General
}

// NewKing Конструктор короля.
func NewKing(pos *entities.Position) *King {
	return &King{
		General: General{
			id:       uuid.New(),
			name:     model.KING,
			Position: pos,
			onBoard:  true,
		},
	}
}

// Rules Проверка хода на соответствие правилам хода фигуры.
func (k *King) Rules(pos *entities.Position) bool {
	return (k.X()-1 == pos.X() || k.X() == pos.X() || k.X()+1 == pos.X()) &&
		(k.Y()-1 == pos.Y() || k.Y() == pos.Y() || k.Y()+1 == pos.Y()) ||
		(!k.AlreadyMove() &&
			(pos.Y() == 1 || pos.Y() == 8) &&
			(pos.X() == 7 || pos.X() == 3))
}

// Направления.
func (k *King) Directions() []entities.Direction {
	return []entities.Direction{
		entities.Top,
		entities.TopRight,
		entities.Right,
		entities.RightBottom,
		entities.Bottom,
		entities.BottomLeft,
		entities.Left,
		entities.LeftTop,
	}
}

// Максимальная дальность хода.
func (k *King) MaxRemote() uint8 {
	return 1
}
