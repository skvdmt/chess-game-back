package chess

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/skvdmt/chess-game-back/internal/entities/dto"
	"github.com/skvdmt/chess-game-back/internal/model"
)

// Clock Часы.
type Clock struct {
	ticker       *time.Ticker
	done         chan struct{}
	step         time.Duration
	reserveWhite time.Duration
	reserveBlack time.Duration
	team         string
}

// NewClock Конструктор.
func NewClock(team string) *Clock {
	return &Clock{
		step:         time.Duration(model.Config.Game.StepTimeLeft) * time.Second,
		reserveWhite: time.Duration(model.Config.Game.ReserveTimeLeft) * time.Second,
		reserveBlack: time.Duration(model.Config.Game.ReserveTimeLeft) * time.Second,
		team:         team,
	}
}

// Start Старт.
func (c *Clock) Start() error {
	model.Logs.Info.Info("board clock starting")
	if c.ticker != nil {
		// Часы уже тикают.
		return nil
	}
	c.done = make(chan struct{}, 1)
	c.ticker = time.NewTicker(time.Second)
	for {
		select {
		case <-c.ticker.C:
			if err := c.tick(); err != nil {
				return err
			}
			j, _ := json.Marshal(c.ExportTick())
			model.CM.Broadcast(j)
		case <-c.done:
			return nil
		}
	}
}

// Stop Остановка.
func (c *Clock) Stop() {
	if c.ticker == nil || c.done == nil {
		return
	}
	c.ticker.Stop()
	c.ticker = nil
	close(c.done)
	c.done = nil
	model.Logs.Info.Info("board clock stopped")
}

// Toggle Переключение.
func (c *Clock) Toggle() {
	c.step = time.Duration(model.Config.Game.StepTimeLeft) * time.Second
	if c.team == model.White {
		c.team = model.Black
		return
	}
	c.team = model.White
}

// Export Експорт.
func (c *Clock) ExportClock() *dto.Clock {
	return &dto.Clock{
		TurnStep:     c.step,
		WhiteReserve: c.reserveWhite,
		BlackReserve: c.reserveBlack,
	}
}

// Export Експорт.
func (c *Clock) ExportTick() *dto.Tick {
	t := &dto.Tick{
		Id:       uuid.New(),
		Name:     "clock",
		TurnStep: c.step,
	}
	switch c.team {
	case model.White:
		t.TurnReserve = c.reserveWhite
	case model.Black:
		t.TurnReserve = c.reserveBlack
	}
	return t
}

// tick Тик.
func (c *Clock) tick() error {
	if c.step > 0 {
		c.step -= time.Second
		return nil
	}
	switch c.team {
	case model.White:
		if c.reserveWhite > 0 {
			c.reserveWhite -= time.Second
			return nil
		}
		return model.ErrTimeOver
	case model.Black:
		if c.reserveBlack > 0 {
			c.reserveBlack -= time.Second
			return nil
		}
		return model.ErrTimeOver
	}
	return nil
}
