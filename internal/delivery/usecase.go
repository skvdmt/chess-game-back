package delivery

import (
	"context"

	"github.com/skvdmt/chess-game-back/internal/entities"
)

// Usecase Интерфейс сервисного слоя.
type Usecase interface {
	// Игра.
	Start(context.Context) error
	// Остановка.
	Stop(ctx context.Context) error
	// Запуск часов.
	StartClock()
	// Остановка часов.
	StopClock()
	// Мат.
	Mate() error
	// Состояние
	State() error
	// Чей ход.
	Turn() string
	// Експорт.
	Export(name string) any
	// Ход.
	Move(from, to *entities.Position) error
	// Новая игра
	New() error
	// Сдаться.
	Surrender(teamName string) error
	// Предложить ничью.
	OfferADraw(teamName string) error
	// Принятие предложения ничьи.
	AcceptDraw(teamName string) error
	// Отклонение предложения ничьи.
	RejectDraw(teamName string) error
}
