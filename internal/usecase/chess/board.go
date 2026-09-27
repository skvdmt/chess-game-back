package chess

import (
	"errors"
	"fmt"

	"github.com/skvdmt/chess-game-back/internal/entities"
	"github.com/skvdmt/chess-game-back/internal/entities/dto"
	"github.com/skvdmt/chess-game-back/internal/model"
)

// Board Доска.
type Board struct {
	// Белые.
	white *Team
	// Черные.
	black *Team
	// Наблюдатели.
	spectators *Team
	// Запущено.
	started bool
	// Чей ход.
	turn string
	// Часы.
	clock *Clock
	// История
	history *History
	// Состояние
	state error
	// Обработчик предложений ничьи.
	offerDraw *OfferDraw
}

// NewBoard Конструктор.
func NewBoard() *Board {
	h := NewHistory()
	b := &Board{
		white:      NewTeam(model.White, h),
		black:      NewTeam(model.Black, h),
		spectators: NewTeam(model.Spectators, nil),
		clock:      NewClock(model.White),
		turn:       model.White,
		history:    h,
		offerDraw:  NewOfferDraw(),
	}
	// Противник.
	b.white.enemy = b.black
	b.black.enemy = b.white
	return b
}

// Start Запуск.
func (b *Board) Start() {
	model.Logs.Info.Info("board starting")
	// Запущено.
	if b.started {
		return
	}
	b.started = true
	b.Turn().CalculatePossibleMoves()
	// Не может сделать ход.
	if !b.Turn().CanMove() {
		b.Pause(b.Turn().Enemy().Win())
		b.Stop()
		return
	}
}

// Stop Остановка.
func (b *Board) Stop() {
	// Остановлено.
	if !b.started {
		return
	}
	b.started = false
	// Остановка чалов.
	b.clockStop()
	// Остановка обработчика предложений ничьи.
	b.offerDraw.Stop()
	model.Logs.Info.Info("board stopped")
}

// Move Ход.
func (b *Board) Move(from, to *entities.Position) error {
	c := b.Turn().ChessPieces().FoundByPosition(from)
	if c == nil {
		return model.ErrChessPieceNotFound
	}
	if !c.Rules(to) {
		switch c.Name() {
		case model.ROOK:
			return model.ErrRulesRook
		case model.KNIGHT:
			return model.ErrRulesKnight
		case model.BISHOP:
			return model.ErrRulesBishop
		case model.QUEEN:
			return model.ErrRulesQueen
		case model.KING:
			return model.ErrRulesKing
		case model.PAWN:
			return model.ErrRulesPawn
		default:
			panic(fmt.Sprintf("unknown chess piece %s", c.Name()))
		}
	}
	if !b.moveValid(c, to) {
		return model.ErrMoveImpossible
	}
	m := NewMove(b.Turn(), c, to, true, b.history)
	m.Do()
	model.CM.SendMove(&dto.Position{
		X: from.X(),
		Y: from.Y(),
	}, &dto.Position{
		X: to.X(),
		Y: to.Y(),
	})
	b.history.Add(m)
	b.toggleTurn()
	return nil
}

// Mate Мат.
func (b *Board) Mate() error {
	b.Turn().CalculatePossibleMoves()
	if !b.Turn().CanMove() {
		if !b.Turn().check() {
			return model.ErrDraw
		}
		switch b.Turn().Name() {
		case model.White:
			return model.ErrBlackWin
		case model.Black:
			return model.ErrWhiteWin
		default:
			panic(fmt.Sprintf("unknown team name %s", b.Turn().Name()))
		}
	}
	return nil
}

// TeamByName Команда по названию.
func (b *Board) TeamByName(name string) *Team {
	switch name {
	case model.White:
		return b.White()
	case model.Black:
		return b.Black()
	case model.Spectators:
		return b.Spectators()
	default:
		panic(fmt.Sprintf("unknown team name %s", name))
	}
}

// OfferDraw Обработчик предложений ничьи.
func (b *Board) OfferDraw() *OfferDraw {
	return b.offerDraw
}

// moveValid Допустимость хода.
func (b *Board) moveValid(c ChessPiece, to *entities.Position) bool {
	for _, p := range c.PossibleMoves() {
		if p.X() == to.X() && p.Y() == to.Y() {
			return true
		}
	}
	return false
}

// toggleTurn Переключлючение хода.
func (b *Board) toggleTurn() {
	// Переключение часов.
	b.Clock().Toggle()
	if b.turn == model.White {
		b.turn = model.Black
	} else {
		b.turn = model.White
	}
	model.CM.SendTurn(b.ExportTurn())
}

// White Белые.
func (b *Board) White() *Team {
	return b.white
}

// Black Черные.
func (b *Board) Black() *Team {
	return b.black
}

// Black Наблюдатели.
func (b *Board) Spectators() *Team {
	return b.spectators
}

// Turn Имя команды чей ход.
func (b *Board) Turn() *Team {
	switch b.turn {
	case model.White:
		return b.white
	case model.Black:
		return b.black
	default:
		panic(fmt.Sprintf("unknown turn %s", b.turn))
	}
}

// Export Экспорт.
func (b *Board) Export() *dto.Board {
	return &dto.Board{
		White: b.white.Export(),
		Black: b.black.Export(),
	}
}

// ExportTurn Экспорт чей ход.
func (b *Board) ExportTurn() *dto.Turn {
	return &dto.Turn{
		Name: b.Turn().Name(),
	}
}

// Clock Часы.
func (b *Board) Clock() *Clock {
	return b.clock
}

// clockStop Остановка часов.
func (b *Board) clockStop() {
	b.clock.Stop()
}

// State Состояние.
func (b *Board) State() error {
	return b.state
}

// pause Приостановление.
func (b *Board) Pause(cause error) {
	if cause == nil {
		panic("cant set pause without cause use unpaused()")
	}
	if errors.Is(b.state, cause) {
		// Причина уже установлена.
		return
	}
	b.state = cause
}
