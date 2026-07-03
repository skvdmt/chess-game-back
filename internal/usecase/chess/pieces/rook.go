package pieces

import (
	"github.com/google/uuid"
	"github.com/skvdmt/chess-game-back/internal/entities"
	"github.com/skvdmt/chess-game-back/internal/model"
)

// Rook Ладья.
type Rook struct {
	General
}

// NewRook Конструктор ладьи.
func NewRook(pos *entities.Position) *Rook {
	return &Rook{
		General: General{
			id:       uuid.New(),
			name:     model.ROOK,
			Position: pos,
			onBoard:  true,
		},
	}
}

// Rules Проверка хода на соответствие правилам хода фигуры.
func (r *Rook) Rules(pos *entities.Position) bool {
	return (r.X() == pos.X() || r.Y() == pos.Y())
}

// Направления.
func (r *Rook) Directions() []entities.Direction {
	return []entities.Direction{
		entities.Top,
		entities.Right,
		entities.Bottom,
		entities.Left,
	}
}

// Максимальная дальность хода.
func (r *Rook) MaxRemote() uint8 {
	return 8
}
