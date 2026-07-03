package pieces

import (
	"github.com/google/uuid"
	"github.com/skvdmt/chess-game-back/internal/entities"
	"github.com/skvdmt/chess-game-back/internal/model"
)

// Queen Королева.
type Queen struct {
	General
}

// NewQueen Конструктор королевы.
func NewQueen(id uuid.UUID, pos *entities.Position) *Queen {
	return &Queen{
		General: General{
			id:       id,
			name:     model.QUEEN,
			Position: pos,
			onBoard:  true,
		},
	}
}

// Rules Проверка хода на соответствие правилам хода фигуры.
func (q *Queen) Rules(pos *entities.Position) bool {
	return (q.X() == pos.X() || q.Y() == pos.Y()) ||
		(pos.X() < q.X() && pos.Y() < q.Y() && q.X()-pos.X() == q.Y()-pos.Y()) ||
		(pos.X() < q.X() && pos.Y() > q.Y() && q.X()-pos.X() == pos.Y()-q.Y()) ||
		(pos.X() > q.X() && pos.Y() < q.Y() && pos.X()-q.X() == q.Y()-pos.Y()) ||
		(pos.X() > q.X() && pos.Y() > q.Y() && pos.X()-q.X() == pos.Y()-q.Y())
}

// Направления.
func (q *Queen) Directions() []entities.Direction {
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
func (q *Queen) MaxRemote() uint8 {
	return 8
}
