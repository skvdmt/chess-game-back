package chess

import (
	"fmt"
	"sync"
	"time"

	"github.com/skvdmt/chess-game-back/internal/model"
)

// OfferDraw Предложение ничьи.
type OfferDraw struct {
	whiteTicker   *time.Ticker
	whiteAttempts int
	whiteLeft     int
	whiteDone     chan struct{}
	blackTicker   *time.Ticker
	blackAttempts int
	blackLeft     int
	blackDone     chan struct{}
	wg            *sync.WaitGroup
}

// NewOfferDraw Конструктор.
func NewOfferDraw() *OfferDraw {
	return &OfferDraw{
		whiteLeft:     model.Config.Game.TimeLeftForConfirmDraw,
		whiteAttempts: model.Config.Game.OfferDrawTimesLeft,
		blackLeft:     model.Config.Game.TimeLeftForConfirmDraw,
		blackAttempts: model.Config.Game.OfferDrawTimesLeft,
		wg:            &sync.WaitGroup{},
	}
}

// Stop Остановка.
func (o *OfferDraw) Stop() {
	if o.whiteDone != nil {
		close(o.whiteDone)
		o.whiteDone = nil
	}
	if o.blackDone != nil {
		close(o.blackDone)
		o.blackDone = nil
	}
	o.wg.Wait()
	model.Logs.Info.Info("offer a draw all handlers stopped")
}

// SendOffer Отправка предложения.
func (o *OfferDraw) SendOffer(t *Team) error {
	if err := o.valid(t); err != nil {
		return err
	}
	o.wg.Add(1)
	go o.tickerHandler(t)
	return nil
}

// Accept Принятие.
func (o *OfferDraw) Accept(t *Team) error {
	if (t.Name() == model.White && o.blackTicker == nil) ||
		(t.Name() == model.Black && o.whiteTicker == nil) {
		return model.ErrNoOfferADraw
	}
	return nil
}

// Reject Отклонение.
func (o *OfferDraw) Reject(t *Team) error {
	model.CM.SendNotice(t.Enemy().Name(), "Draw Rejected")
	if t.Name() == model.White && o.blackTicker != nil {
		close(o.blackDone)
		o.blackDone = nil
		return nil
	}
	if t.Name() == model.Black && o.whiteTicker != nil {
		close(o.whiteDone)
		o.whiteDone = nil
		return nil
	}
	return nil
}

// valid Валидация.
func (o *OfferDraw) valid(t *Team) error {
	if t.Name() == model.Spectators {
		return model.ErrYouAreSpectator
	}
	if (t.Name() == model.White && o.whiteAttempts == 0) ||
		(t.Name() == model.Black && o.blackAttempts == 0) {
		return model.ErrAttemptsOfferADrawOver
	}
	if (t.Name() == model.White && o.whiteTicker != nil) ||
		(t.Name() == model.Black && o.blackTicker != nil) {
		return model.ErrOfferADrawAlreadySent
	}
	return nil
}

// tickerHandler Обработчик тикера.
func (o *OfferDraw) tickerHandler(t *Team) {
	model.Logs.Info.Info(fmt.Sprintf("offer a draw %s starting", t.Name()))
	switch t.Name() {
	case model.White:
		o.whiteTicker = time.NewTicker(time.Second)
		o.whiteDone = make(chan struct{}, 1)
		defer func() {
			o.whiteAttempts--
			o.whiteLeft = model.Config.Game.TimeLeftForConfirmDraw
			o.whiteTicker.Stop()
			o.whiteTicker = nil
			model.Logs.Info.Info(fmt.Sprintf("offer a draw %s stopped", t.Name()))
			o.wg.Done()
		}()
		for {
			select {
			case <-o.whiteTicker.C:
				o.whiteLeft--
				model.CM.SendOfferADraw(
					t.Enemy().Name(),
					time.Duration(o.whiteLeft),
				)
				if o.whiteLeft <= 0 {
					model.CM.SendNotice(t.Name(), "Draw Rejected")
					return
				}
			case <-o.whiteDone:
				return
			}
		}
	case model.Black:
		o.blackTicker = time.NewTicker(time.Second)
		o.blackDone = make(chan struct{}, 1)
		defer func() {
			o.blackAttempts--
			o.blackLeft = model.Config.Game.TimeLeftForConfirmDraw
			o.blackTicker.Stop()
			o.blackTicker = nil
			model.Logs.Info.Info(fmt.Sprintf("offer a draw %s stopped", t.Name()))
			o.wg.Done()
		}()
		for {
			select {
			case <-o.blackTicker.C:
				o.blackLeft--
				model.CM.SendOfferADraw(
					t.Enemy().Name(),
					time.Duration(o.blackLeft),
				)
				if o.blackLeft <= 0 {
					model.CM.SendNotice(t.Name(), "Draw Rejected")
					return
				}
			case <-o.blackDone:
				return
			}
		}
	}
}
