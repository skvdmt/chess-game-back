package chess

import (
	"github.com/skvdmt/chess-game-back/internal/entities"
	"github.com/skvdmt/chess-game-back/internal/model"
)

// History История ходов.
type History []*Move

// NewHistory Конструктор.
func NewHistory() *History {
	return &History{}
}

// PassantCapture Допустимо взятие на проходе.
func (h *History) PassantCapture(t *Team, c ChessPiece, p *entities.Position) bool {
	if len(*h) > 0 {
		pm := (*h)[len(*h)-1]
		if c.Name() == model.PAWN &&
			t.Enemy().ChessPieces().FoundById(pm.primaryId).Name() == model.PAWN &&
			pm.primaryStart.X() == p.X() &&
			((p.Y() == 6 && pm.primaryStart.Y() == 7 &&
				pm.primaryFinish.Y() == 5) ||
				(p.Y() == 3 && pm.primaryStart.Y() == 2 &&
					pm.primaryFinish.Y() == 4)) {
			// Взятие на проходе.
			return true
		}
	}
	return false
}

// Add Добавить.
func (h *History) Add(m *Move) {
	*h = append(*h, m)
}
