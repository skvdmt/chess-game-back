package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ChessPiece Шахматная фигура.
type ChessPiece struct {
	Name     string    `json:"name"`
	Position *Position `json:"position"`
}

// Часы.
type Clock struct {
	TurnStep     time.Duration `json:"turn_step"`
	WhiteReserve time.Duration `json:"white_reserve"`
	BlackReserve time.Duration `json:"black_reserve"`
}

// Tick Тик событие.
type Tick struct {
	Id          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	TurnStep    time.Duration `json:"turn_step"`
	TurnReserve time.Duration `json:"turn_reserve"`
}

// Board Доска.
type Board struct {
	White *Team `json:"white"`
	Black *Team `json:"black"`
}

// Event Событие.
type Event struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Body any       `json:"body"`
}

// Ход.
type Move struct {
	From *Position `json:"from"`
	To   *Position `json:"to"`
}

// Position Позиция.
type Position struct {
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}

// Request Запрос.
type Request struct {
	Id     uuid.UUID       `json:"id"`
	Method string          `json:"method"`
	Body   json.RawMessage `json:"body"`
}

// Response Ответ.
type Response struct {
	Id        uuid.UUID `json:"id"`
	RequestId uuid.UUID `json:"request_id"`
	Body      any       `json:"body"`
}

// Статус.
type Status struct {
	Valid bool   `json:"valid"`
	Cause string `json:"cause,omitempty"`
}

// TeamName Имя команды.
type TeamName struct {
	Name string `json:"name"`
}

// Team Команда.
type Team struct {
	Name               string        `json:"name"`
	OnBoardChessPieces []*ChessPiece `json:"on_board_chess_pieces"`
	EatenChessPieces   []*ChessPiece `json:"eaten_chess_pieces"`
}

// Turn Чей ход.
type Turn struct {
	Name string `json:"name"`
}

// OfferADraw Предложить ничью.
type OfferADraw struct {
	DesisionTimeLeft time.Duration `json:"desision_time_left"`
}

// Notice Нотация.
type Notice struct {
	Notice string `json:"notice"`
}
