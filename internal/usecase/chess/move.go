package chess

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/skvdmt/chess-game-back/internal/entities"
	"github.com/skvdmt/chess-game-back/internal/entities/dto"
	"github.com/skvdmt/chess-game-back/internal/model"
	"github.com/skvdmt/chess-game-back/internal/usecase/chess/pieces"
)

const (
	// Ход.
	General MoveType = iota + 1
	// Трансформация.
	Transform
	// Поедание.
	Eating
	// Поедание с трансформацией.
	TransformEating
	// Рокировка.
	Castling
	// Взятие на проходе.
	PassantCapture
)

// MoveType Тип хода.
type MoveType uint8

// Move Ход.
type Move struct {
	// Id хода.
	id uuid.UUID
	// Команда.
	team *Team
	// Фигура уже ходила.
	alreadyMove bool
	// Тип.
	moveType MoveType
	// Id основной фигуры.
	primaryId uuid.UUID
	// Id дополнительной фигуры.
	secondaryId uuid.UUID
	// Стартовая позиция основной фигуры.
	primaryStart *entities.Position
	// Конечная позиция основной фигуры.
	primaryFinish *entities.Position
	// Стартовая позиция дополнительной фигуры.
	secondaryStart *entities.Position
	// Конечная позиция дополнительной фигуры.
	secondaryFinish *entities.Position
	// Фигура из которой произмодится трансформация.
	prev ChessPiece
	// Реальный ход.
	real bool
}

// NewMove Конструктор.
func NewMove(t *Team, c ChessPiece, p *entities.Position, real bool, h *History) *Move {
	m := &Move{
		id:            uuid.New(),
		team:          t,
		alreadyMove:   c.AlreadyMove(),
		primaryId:     c.Id(),
		primaryStart:  c.Pos(),
		primaryFinish: p,
		real:          real,
	}
	m.moveType = m.typ(p, h)

	switch m.moveType {
	case General, Transform:
	case Eating, TransformEating:
		m.secondaryId = m.team.Enemy().ChessPieces().FoundByPosition(p).Id()
	case PassantCapture:
		switch t.Name() {
		case model.White:
			m.secondaryId = m.team.Enemy().ChessPieces().FoundByPosition(entities.NewPosition(p.X(), p.Y()-1)).Id()
		case model.Black:
			m.secondaryId = m.team.Enemy().ChessPieces().FoundByPosition(entities.NewPosition(p.X(), p.Y()+1)).Id()
		}
	case Castling:
		switch p.X() {
		case 3:
			r := m.team.ChessPieces().FoundByPosition(entities.NewPosition(1, c.Pos().Y()))
			m.secondaryId = r.Id()
			m.secondaryStart = r.Pos()
			m.secondaryFinish = entities.NewPosition(4, r.Pos().Y())
		case 7:
			r := m.team.ChessPieces().FoundByPosition(entities.NewPosition(8, c.Pos().Y()))
			m.secondaryId = r.Id()
			m.secondaryStart = r.Pos()
			m.secondaryFinish = entities.NewPosition(6, r.Pos().Y())
		}
	}
	return m
}

// Do Выполнение.
func (m *Move) Do() {
	c := m.team.ChessPieces().FoundById(m.primaryId)
	c.SetPos(m.primaryFinish, m.real)
	switch m.moveType {
	case General:
		break
	case Transform:
		m.prev = c
		q := pieces.NewQueen(m.primaryId, m.primaryFinish)
		m.team.ChessPieces().Replace(c, q)
	case Eating:
		e := m.team.Enemy().ChessPieces().FoundById(m.secondaryId)
		e.Eat()
	case TransformEating:
		e := m.team.Enemy().ChessPieces().FoundById(m.secondaryId)
		e.Eat()
		m.prev = c
		q := pieces.NewQueen(m.primaryId, m.primaryFinish)
		m.team.ChessPieces().Replace(c, q)
	case PassantCapture:
		e := m.team.Enemy().ChessPieces().FoundById(m.secondaryId)
		e.Eat()
	case Castling:
		s := m.team.ChessPieces().FoundById(m.secondaryId)
		model.CM.SendMove(&dto.Position{
			X: s.X(),
			Y: s.Y(),
		},
			&dto.Position{
				X: m.secondaryFinish.X(),
				Y: m.secondaryFinish.Y(),
			})
		s.SetPos(m.secondaryFinish, m.real)
	default:
		panic(fmt.Sprintf("unknown move type %d", m.moveType))
	}
}

// Undo Отмена.
func (m *Move) Undo() {
	c := m.team.ChessPieces().FoundById(m.primaryId)
	switch m.moveType {
	case General:
		break
	case Transform:
		m.team.ChessPieces().Replace(c, m.prev)
		c = m.team.ChessPieces().FoundById(m.primaryId)
	case Eating:
		e := m.team.Enemy().ChessPieces().FoundById(m.secondaryId)
		e.UndoEat()
	case PassantCapture:
		e := m.team.Enemy().ChessPieces().FoundById(m.secondaryId)
		e.UndoEat()
	case TransformEating:
		e := m.team.Enemy().ChessPieces().FoundById(m.secondaryId)
		e.UndoEat()
		m.team.ChessPieces().Replace(c, m.prev)
		c = m.team.ChessPieces().FoundById(m.primaryId)
	case Castling:
		s := m.team.ChessPieces().FoundById(m.secondaryId)
		s.SetPos(m.secondaryStart, m.real)
	default:
		panic(fmt.Sprintf("unknown move type %d", m.moveType))
	}
	c.SetPos(m.primaryStart, m.real)
}

// typ Тип.
func (m *Move) typ(p *entities.Position, h *History) MoveType {
	c := m.team.ChessPieces().FoundById(m.primaryId)
	if c.Name() == model.KING &&
		!c.AlreadyMove() &&
		(p.X() == 3 || p.X() == 7) &&
		(p.Y() == 1 || p.Y() == 8) {
		// Рокировка
		return Castling
	}
	if m.team.Enemy().ChessPieces().FoundByPosition(p) != nil {
		if c.Name() == model.PAWN &&
			(p.Y() == 1 || p.Y() == 8) {
			// Поедание с трансформацией.
			return TransformEating
		}
		// Поедание.
		return Eating
	}
	if c.Name() == model.PAWN &&
		(p.Y() == 1 || p.Y() == 8) {
		// Трансформация.
		return Transform
	}
	if len(*h) > 0 {
		pm := (*h)[len(*h)-1]
		if c.Name() == model.PAWN &&
			m.team.Enemy().ChessPieces().FoundById(pm.primaryId).Name() == model.PAWN &&
			pm.primaryStart.X() == p.X() &&
			((p.Y() == 6 && pm.primaryStart.Y() == 7 &&
				pm.primaryFinish.Y() == 5) ||
				(p.Y() == 3 && pm.primaryStart.Y() == 2 &&
					pm.primaryFinish.Y() == 4)) {
			// Взятие на проходе.
			return PassantCapture
		}
	}
	return General
}
