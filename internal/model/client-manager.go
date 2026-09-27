package model

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/skvdmt/chess-game-back/internal/entities/dto"
)

const (
	White      = "white"
	Black      = "black"
	Spectators = "spectators"
)

// CM Глобальный менеджер клиентов.
var CM *ClientManager

// ClientManager Управление клиентами.
type ClientManager struct {
	// Клиенты.
	clients map[*Client]bool
	// Обработчики.
	handlers map[string]func(Client, *dto.Request) any
	// Мьютекс.
	mu *sync.Mutex
	// Последний отправленый статус.
	lastStatus *dto.Status
}

// CreateClientManager Конструктор
func CreateClientManager() {
	CM = &ClientManager{
		clients:    make(map[*Client]bool),
		handlers:   make(map[string]func(Client, *dto.Request) any),
		mu:         &sync.Mutex{},
		lastStatus: &dto.Status{},
	}
}

// Close Закрытие.
func (c *ClientManager) Close() {
	// TODO удалить всех клиентов.
}

// Register Регистрация клиента.
func (c *ClientManager) Register(client *Client) {
	client.manager = c
	c.mu.Lock()
	c.clients[client] = true
	c.mu.Unlock()
}

// Unregister Анрегистрация клиента.
func (c *ClientManager) Unregister(client *Client) {
	c.mu.Lock()
	delete(c.clients, client)
	c.mu.Unlock()
	Logs.Info.Info(fmt.Sprintf("client %s disconnected", client.teamName))
}

// Broadcast Отправка сообщения всем клиентам.
func (c *ClientManager) Broadcast(message []byte) {
	for client := range c.clients {
		select {
		case client.send <- message:
		default:
			client.Close()
			c.Unregister(client)
		}
	}
}

// Join Присоединение клиента к команде.
func (c *ClientManager) Join(client *Client) {
	if !c.ClientTeamExists(White) {
		Logs.Info.Info(fmt.Sprintf("client connected as %s", White))
		client.teamName = White
		return
	}
	if !c.ClientTeamExists(Black) {
		Logs.Info.Info(fmt.Sprintf("client connected as %s", Black))
		client.teamName = Black
		return
	}
	Logs.Info.Info(fmt.Sprintf("client connected as %s", Spectators))
	client.teamName = Spectators
}

// ClientTeamExists Клиент команды существует.
func (c *ClientManager) ClientTeamExists(teamName string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for c := range c.clients {
		if c.teamName == teamName {
			return true
		}
	}
	return false
}

// Handle Регистрация обработчика.
func (c *ClientManager) Handle(method string, handler func(Client, *dto.Request) any) {
	c.handlers[method] = handler
}

// SendStatus Отправка статуса.
func (c *ClientManager) SendStatus(valid bool, cause string) {
	if c.lastStatus.Valid == valid && c.lastStatus.Cause == cause {
		// Текущий статус уже был отправлен
		return
	}
	j, _ := json.Marshal(&dto.Event{
		Id:   uuid.New(),
		Name: "status",
		Body: &dto.Status{
			Valid: valid,
			Cause: cause,
		},
	})
	c.Broadcast(j)
	c.lastStatus = &dto.Status{
		Valid: valid,
		Cause: cause,
	}
}

// SendMove Отправка хода.
func (c *ClientManager) SendMove(from, to *dto.Position) {
	j, _ := json.Marshal(&dto.Event{
		Id:   uuid.New(),
		Name: "move",
		Body: &dto.Move{
			From: from,
			To:   to,
		},
	})
	c.Broadcast(j)
}

// SendMove Отправка хода.
func (c *ClientManager) SendTurn(t *dto.Turn) {
	j, _ := json.Marshal(&dto.Event{
		Id:   uuid.New(),
		Name: "turn",
		Body: t,
	})
	c.Broadcast(j)
}

// Swap Смена команд.
func (c *ClientManager) Swap() {
	var wc, bc *Client
	for c := range c.clients {
		if c.teamName == White {
			wc = c
		}
		if c.teamName == Black {
			bc = c
		}
	}
	if wc != nil {
		wc.teamName = Black
	}
	if bc != nil {
		bc.teamName = White
	}
}

// SendOfferADraw Отправить предложение зафиксировать ничью противнику.
func (c *ClientManager) SendOfferADraw(teamName string, left time.Duration) {
	j, _ := json.Marshal(&dto.Event{
		Id:   uuid.New(),
		Name: "offer_a_draw",
		Body: &dto.OfferADraw{
			DesisionTimeLeft: left,
		},
	})
	for c := range c.clients {
		if c.teamName == teamName {
			c.Send(j)
			return
		}
	}
}

// SendNotice Отправка нотации.
func (c *ClientManager) SendNotice(teamName, notice string) {
	j, _ := json.Marshal(&dto.Event{
		Id:   uuid.New(),
		Name: "notice",
		Body: &dto.Notice{
			Notice: notice,
		},
	})
	for c := range c.clients {
		if c.teamName == teamName {
			c.Send(j)
			return
		}
	}
}
