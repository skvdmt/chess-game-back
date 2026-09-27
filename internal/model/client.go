package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/skvdmt/chess-game-back/internal/entities/dto"
)

const (
	maxMessageSize = 4096
	pongWait       = time.Second * 15
	pingPeriod     = (pongWait * 9) / 10
	writeWait      = time.Second * 15

	// Предоставить всю информацию о доске
	RequestMethodGetBoard       = "GET_BOARD"
	RequestMethodGetTurn        = "GET_TURN"
	RequestMethodGetTeamName    = "GET_TEAM_NAME"
	RequestMethodGetStatus      = "GET_STATUS"
	RequestMethodGetClock       = "GET_CLOCK"
	RequestMethodPostMove       = "POST_MOVE"
	RequestMethodPostNewGame    = "POST_NEW_GAME"
	RequestMethodPostSurrender  = "POST_SURRENDER"
	RequestMethodPostOfferADraw = "POST_OFFER_A_DRAW"
	RequestMethodPostAcceptDraw = "POST_ACCEPT_DRAW"
	RequestMethodPostRejectDraw = "POST_REJECT_DRAW"
)

// Client Клиент для соединения с сервером по WebSocket.
type Client struct {
	manager  *ClientManager
	con      *websocket.Conn
	teamName string
	send     chan []byte
}

// NewClient Конструктор.
func NewClient(con *websocket.Conn) *Client {
	c := &Client{
		con:  con,
		send: make(chan []byte, 256),
	}
	go c.read()
	go c.write()
	return c
}

// TeamName Имя команды.
func (c *Client) TeamName() string {
	return c.teamName
}

// Send Отправка сообщения клиенту.
func (c *Client) Send(msg []byte) {
	c.send <- msg
}

// Close Закрытие клиента.
func (c *Client) Close() {
	close(c.send)
}

var (
	ErrBadRequest = errors.New("bad request")
)

// Read Прием данных клиентом от сервера.
func (c *Client) read() {
	defer func() {
		c.manager.Unregister(c)
		_ = c.con.Close()
	}()
	c.con.SetReadLimit(maxMessageSize)
	_ = c.con.SetReadDeadline(time.Now().Add(pongWait))
	c.con.SetPongHandler(func(string) error {
		_ = c.con.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		var req *dto.Request
		if err := c.con.ReadJSON(&req); err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				break
			}
			Logs.Error.Error(err.Error())
			break
		}
		c.requestHandler(req)
	}
}

// Write Отправка данных клиентом серверу.
func (c *Client) write() {
	ping := time.NewTicker(pingPeriod)
	defer func() {
		ping.Stop()
		_ = c.con.Close()
	}()
	for {
		select {
		case m := <-c.send:
			_ = c.con.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.con.WriteMessage(websocket.TextMessage, m); err != nil {
				Logs.Info.Info(err.Error())
			}
		case <-ping.C:
			_ = c.con.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.con.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// requestHandler Обработка запроса.
func (c *Client) requestHandler(req *dto.Request) {
	for m, h := range c.manager.handlers {
		if req.Method == m {
			res := &dto.Response{
				Id:        uuid.New(),
				RequestId: req.Id,
			}
			if req.Method == m {
				res.Body = h(*c, req)
			}
			c.sendResponse(res)
			return
		}
	}
	c.sendResponse(&dto.Response{
		Id:        uuid.New(),
		RequestId: req.Id,
		Body: &dto.Status{
			Valid: false,
			Cause: "unknown method",
		},
	})
}

// sendResponse Отправка ответа на запрос.
func (c *Client) sendResponse(res *dto.Response) {
	j, err := json.Marshal(res)
	if err != nil {
		Logs.Error.Error(err.Error())
		return
	}
	c.Send(j)
}
