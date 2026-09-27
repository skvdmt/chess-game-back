package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/skvdmt/chess-game-back/internal/entities"
	"github.com/skvdmt/chess-game-back/internal/entities/dto"
	"github.com/skvdmt/chess-game-back/internal/model"
	"github.com/skvdmt/chess-game-back/internal/usecase/chess"
)

// App Сервисный слой.
type App struct {
	// Доска.
	board *chess.Board
}

// NewApp Конструктор.
func NewApp(_ context.Context) (*App, error) {
	model.Logs.Info.Info("usecase layer creating")
	a := &App{
		board: chess.NewBoard(),
	}
	return a, nil
}

// Start Запуск.
func (a *App) Start(_ context.Context) error {
	model.Logs.Info.Info("usecase layer starting")
	// Запуск доски
	a.board.Start()
	return nil
}

// Stop Остановка.
func (a *App) Stop(_ context.Context) error {
	a.board.Stop()
	model.Logs.Info.Info("usecase layer stopped")
	return nil
}

// StartClock Запуск часов.
func (a *App) StartClock() {
	// Запуск часов.
	if err := a.board.Clock().Start(); err != nil {
		if errors.Is(err, model.ErrTimeOver) {
			a.board.Pause(a.board.Turn().Enemy().Win())
		}
		a.board.Stop()
		return
	}
	a.board.Clock().Stop()
}

// StopClock Остановка часов.
func (a *App) StopClock() {
	a.board.Clock().Stop()
}

// Turn Чей ход.
func (a *App) Turn() string {
	return a.board.Turn().Name()
}

// Export Експортирование.
func (a *App) Export(name string) any {
	switch name {
	case "board":
		return a.board.Export()
	case "turn":
		return a.board.ExportTurn()
	case "clock":
		return a.board.Clock().ExportClock()
	default:
		panic(fmt.Sprintf("unknown name %s", name))
	}
}

// Move Сделать ход.
func (a *App) Move(from, to *entities.Position) error {
	// Проверка сервисного слоя на наличие паузы.
	if err := a.State(); err != nil {
		return err
	}
	if err := a.board.Move(from, to); err != nil {
		return err
	}
	return nil
}

// Mate Мат.
func (a *App) Mate() error {
	if err := a.board.Mate(); err != nil {
		a.board.Pause(err)
		a.board.Stop()
		return err
	}
	return nil
}

// New Новая игра.
func (a *App) New() error {
	if !errors.Is(a.State(), model.ErrWhiteWin) &&
		!errors.Is(a.State(), model.ErrBlackWin) &&
		!errors.Is(a.State(), model.ErrDraw) {
		return model.ErrGameNotOver
	}
	if model.Config.Game.SwapTeamsAfterMakingNewGame {
		model.CM.Swap()
	}
	a.board.Stop()
	a.board = chess.NewBoard()
	// Запуск доски
	a.board.Start()
	j, _ := json.Marshal(&dto.Event{
		Id:   uuid.New(),
		Name: "created",
		Body: nil,
	})
	model.CM.Broadcast(j)
	return nil
}

// Surrender Сдаться.
func (a *App) Surrender(teamName string) error {
	var err error
	switch teamName {
	case model.White:
		err = model.ErrBlackWin
	case model.Black:
		err = model.ErrWhiteWin
	default:
		panic(fmt.Sprintf("unknown team name %s", teamName))
	}
	a.board.Pause(err)
	a.board.Stop()
	return nil
}

// Status Статус.
func (a *App) State() error {
	return a.board.State()
}

// OfferADraw Предложить ничью.
func (a *App) OfferADraw(teamName string) error {
	return a.board.OfferDraw().SendOffer(a.board.TeamByName(teamName))
}

// AcceptDraw Принятие предложения ничьи.
func (a *App) AcceptDraw(teamName string) error {
	if err := a.board.OfferDraw().Accept(a.board.TeamByName(teamName)); err != nil {
		return err
	}
	a.board.Pause(model.ErrDraw)
	a.board.Stop()
	return nil
}

// RejectDraw Отклонение предложения ничьи.
func (a *App) RejectDraw(teamName string) error {
	return a.board.OfferDraw().Reject(a.board.TeamByName(teamName))
}
