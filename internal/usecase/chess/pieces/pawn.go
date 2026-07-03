package pieces

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/skvdmt/chess-game-back/internal/entities"
	"github.com/skvdmt/chess-game-back/internal/model"
)

// Pawn Пешка.
type Pawn struct {
	team string
	General
}

// NewPawn Конструктор пешки.
func NewPawn(id uuid.UUID, pos *entities.Position, team string) *Pawn {
	return &Pawn{
		team: team,
		General: General{
			id:       id,
			name:     model.PAWN,
			Position: pos,
			onBoard:  true,
		},
	}
}

// Rules Проверка хода на соответствие правилам хода фигуры.
func (p *Pawn) Rules(pos *entities.Position) bool {
	switch p.team {
	case model.White:
		return (p.X() == pos.X() && p.Y()+1 == pos.Y()) ||
			(p.X() == pos.X()+1 && p.Y()+1 == pos.Y()) ||
			(p.X() == pos.X()-1 && p.Y()+1 == pos.Y()) ||
			(p.X() == pos.X() && p.Y()+2 == pos.Y() && !p.AlreadyMove())
	case model.Black:
		return (p.X() == pos.X() && p.Y()-1 == pos.Y()) ||
			(p.X() == pos.X()+1 && p.Y()-1 == pos.Y()) ||
			(p.X() == pos.X()-1 && p.Y()-1 == pos.Y()) ||
			(p.X() == pos.X() && p.Y()-2 == pos.Y() && !p.AlreadyMove())
	default:
		panic(fmt.Sprintf("unknown team name %s", p.team))
	}
}

// Направления.
func (p *Pawn) Directions() []entities.Direction {
	switch p.team {
	case model.White:
		return []entities.Direction{entities.Top}
	case model.Black:
		return []entities.Direction{entities.Bottom}
	default:
		panic(fmt.Sprintf("unknown team name %s", p.team))
	}
}

// Максимальная дальность хода.
func (p *Pawn) MaxRemote() uint8 {
	if p.AlreadyMove() {
		return 1
	}
	return 2
}
