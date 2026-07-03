package chess

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/skvdmt/chess-game-back/internal/entities"
	"github.com/skvdmt/chess-game-back/internal/entities/dto"
	"github.com/skvdmt/chess-game-back/internal/model"
	"github.com/skvdmt/chess-game-back/internal/usecase/chess/pieces"
)

// Шахматные фигуры.
type ChessPieces []ChessPiece

// NewChessPieces Конструктор.
func NewChessPieces(teamName string) *ChessPieces {
	c := &ChessPieces{}
	var fl, pl uint8
	switch teamName {
	case model.White:
		fl = 1
		pl = 2
	case model.Black:
		fl = 8
		pl = 7
	case model.Spectators:
		return nil
	default:
		panic(fmt.Errorf("unknown team name %s", teamName))
	}
	for x := uint8(1); x <= 8; x++ {
		c.Add(pieces.NewPawn(uuid.New(), entities.NewPosition(x, pl), teamName))
	}
	c.Add(pieces.NewRook(entities.NewPosition(1, fl)))
	c.Add(pieces.NewRook(entities.NewPosition(8, fl)))
	c.Add(pieces.NewKnight(entities.NewPosition(2, fl)))
	c.Add(pieces.NewKnight(entities.NewPosition(7, fl)))
	c.Add(pieces.NewBishop(entities.NewPosition(3, fl)))
	c.Add(pieces.NewBishop(entities.NewPosition(6, fl)))
	c.Add(pieces.NewKing(entities.NewPosition(5, fl)))
	c.Add(pieces.NewQueen(uuid.New(), entities.NewPosition(4, fl)))
	return c
}

// King Получить короля.
func (cs *ChessPieces) King() ChessPiece {
	for _, c := range *cs {
		if c.Name() == model.KING {
			return c
		}
	}
	panic("king not found")
}

// Add Добавление.
func (cs *ChessPieces) Add(c ChessPiece) {
	*cs = append(*cs, c)
}

// Replace Замена.
func (cs *ChessPieces) Replace(o, n ChessPiece) {
	for i, c := range *cs {
		if c.Id() == o.Id() {
			(*cs)[i] = n
			return
		}
	}
	panic(fmt.Sprintf("chess piece with id %s not found", o.Id()))
}

// FoundByPosition Поиск по позиции.
func (cs *ChessPieces) FoundByPosition(p *entities.Position) ChessPiece {
	for _, c := range *cs {
		if !c.OnBoard() {
			continue
		}
		if c.X() == p.X() && c.Y() == p.Y() {
			return c
		}
	}
	return nil
}

// FoundById Поиск по id.
func (cs *ChessPieces) FoundById(id uuid.UUID) ChessPiece {
	for _, c := range *cs {
		if c.Id() == id {
			return c
		}
	}
	return nil
}

// FoundBy Наличие по позиции.
func (cs *ChessPieces) ExistsByPosition(p *entities.Position) bool {
	return cs.FoundByPosition(p) != nil
}

// ExportOnBoadr Экспорт фигур на доске.
func (cs *ChessPieces) ExportOnBoadr() []*dto.ChessPiece {
	return cs.export(true)
}

// ExportEaten Экспорт съеденых фигур.
func (cs *ChessPieces) ExportEaten() []*dto.ChessPiece {
	return cs.export(false)
}

// export Экспорт фигур.
func (cs *ChessPieces) export(onBoard bool) []*dto.ChessPiece {
	var e []*dto.ChessPiece
	for _, c := range *cs {
		if c.OnBoard() != onBoard {
			continue
		}
		e = append(e, &dto.ChessPiece{
			Name: c.Name(),
			Position: &dto.Position{
				X: c.X(),
				Y: c.Y(),
			},
		})
	}
	return e
}
