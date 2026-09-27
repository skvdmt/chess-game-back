package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/skvdmt/chess-game-back/internal/entities"
	"github.com/skvdmt/chess-game-back/internal/entities/dto"
	"github.com/skvdmt/chess-game-back/internal/model"
	"github.com/skvdmt/chess-game-back/internal/usecase"
)

// App Транспортный слой.
type App struct {
	// Триггер остановки.
	stopTrigger chan struct{}
	// Часы запущены.
	clockStarted bool
	// Состояние.
	state error // Ждем обоих | ждем белого | ждем черного
	// Сервисный слой.
	usecase Usecase
	upg     *websocket.Upgrader
	wg      *sync.WaitGroup
}

// NewApp Конструктор.
func NewApp(ctx context.Context) (*App, error) {
	model.Logs.Info.Info("delivery layer creating")
	a := &App{
		stopTrigger: make(chan struct{}, 1),
		wg:          &sync.WaitGroup{},
	}
	// Создание менеджера клиентов.
	model.CreateClientManager()
	a.routes()
	var err error
	// Создание сервисного слоя.
	a.usecase, err = usecase.NewApp(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Start Запуск.
func (a *App) Start(ctx context.Context) error {
	model.Logs.Info.Info("delivery layer starting")
	a.upg = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			if r.Header.Get("origin") != model.Config.Server.OriginalClientUrl {
				return false
			}
			return true
		},
	}
	if err := a.usecase.Start(ctx); err != nil {
		model.Errors <- err
	}
	a.wg.Add(1)
	go a.play()
	return nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	model.Logs.Info.Info("delivery play stopped")
	a.stopTrigger <- struct{}{}
	close(a.stopTrigger)
	a.wg.Wait()
	if err := a.usecase.Stop(ctx); err != nil {
		return err
	}
	model.Logs.Info.Info("delivery layer stopped")
	return nil
}

// Connect Присоединение.
func (a *App) Connect(w http.ResponseWriter, r *http.Request) {
	con, err := a.upg.Upgrade(w, r, nil)
	if err != nil {
		model.Logs.Info.Info(err.Error())
		return
	}
	c := model.NewClient(con)
	model.CM.Register(c)
	model.CM.Join(c)
}

// Event Событие.
func (a *App) Event(name string, b any) {
	j, err := json.Marshal(&dto.Event{
		Id:   uuid.New(),
		Name: name,
		Body: b,
	})
	if err != nil {
		model.Logs.Error.Error(err.Error())
		return
	}
	model.CM.Broadcast(j)
}

// State Состояние.
func (a *App) State() error {
	if a.state != nil {
		return a.state
	}
	return a.usecase.State()
}

// play Игра.
func (a *App) play() {
	defer a.wg.Done()
	model.Logs.Info.Info("delivery play started")
	for {
		time.Sleep(time.Millisecond * 300)
		select {
		case <-a.stopTrigger:
			return
		default:
			a.clientsConnected()
			if err := a.State(); err != nil {
				model.CM.SendStatus(false, err.Error())
				if a.clockStarted {
					a.clockStarted = false
					a.usecase.StopClock()
				}
				// Игра не может начаться.
				continue
			}
			model.CM.SendStatus(true, "")
			if !a.clockStarted {
				a.clockStarted = true
				go a.usecase.StartClock()
			}
		}
	}
}

// clientsConnected Проверка наличия клиентов для возобновления игры.
func (a *App) clientsConnected() {
	we := model.CM.ClientTeamExists(model.White)
	be := model.CM.ClientTeamExists(model.Black)
	// fmt.Println(we, be, a.state)
	switch {
	case !we && !be:
		a.pause(model.ErrWaitBothPlayers)
	case !we:
		a.pause(model.ErrWaitWhitePlayer)
	case !be:
		a.pause(model.ErrWaitBlackPlayer)
	default:
		a.unpause()
	}
}

// pause Приостановление.
func (a *App) pause(cause error) {
	if cause == nil {
		panic("cant set pause without cause use unpaused()")
	}
	if errors.Is(a.state, cause) {
		// Причина уже установлена.
		return
	}
	a.state = cause
}

// unpause Возобновление.
func (a *App) unpause() {
	a.state = nil
}

// routes Маршруты.
func (a *App) routes() {
	model.CM.Handle(model.RequestMethodPostAcceptDraw, func(c model.Client, r *dto.Request) any {
		// Вы наблюдатель.
		if c.TeamName() == model.Spectators {
			return &dto.Status{
				Valid: false,
				Cause: model.ErrYouAreSpectator.Error(),
			}
		}
		// Игра не идет.
		if err := a.State(); err != nil {
			return &dto.Status{
				Valid: false,
				Cause: err.Error(),
			}
		}
		if err := a.usecase.AcceptDraw(c.TeamName()); err != nil {
			return &dto.Status{
				Valid: false,
				Cause: err.Error(),
			}
		}
		return &dto.Status{
			Valid: true,
		}
	})
	model.CM.Handle(model.RequestMethodPostRejectDraw, func(c model.Client, r *dto.Request) any {
		// Вы наблюдатель.
		if c.TeamName() == model.Spectators {
			return &dto.Status{
				Valid: false,
				Cause: model.ErrYouAreSpectator.Error(),
			}
		}
		// Игра не идет.
		if err := a.State(); err != nil {
			return &dto.Status{
				Valid: false,
				Cause: err.Error(),
			}
		}
		if err := a.usecase.RejectDraw(c.TeamName()); err != nil {
			return &dto.Status{
				Valid: false,
				Cause: err.Error(),
			}
		}
		return &dto.Status{
			Valid: true,
		}
	})
	model.CM.Handle(model.RequestMethodPostOfferADraw, func(c model.Client, r *dto.Request) any {
		// Вы наблюдатель.
		if c.TeamName() == model.Spectators {
			return &dto.Status{
				Valid: false,
				Cause: model.ErrYouAreSpectator.Error(),
			}
		}
		// Игра не идет.
		if err := a.State(); err != nil {
			return &dto.Status{
				Valid: false,
				Cause: err.Error(),
			}
		}
		if err := a.usecase.OfferADraw(c.TeamName()); err != nil {
			return &dto.Status{
				Valid: false,
				Cause: err.Error(),
			}
		}
		return &dto.Status{
			Valid: true,
		}
	})

	model.CM.Handle(model.RequestMethodPostSurrender, func(c model.Client, r *dto.Request) any {
		// Вы наблюдатель.
		if c.TeamName() == model.Spectators {
			return &dto.Status{
				Valid: false,
				Cause: model.ErrYouAreSpectator.Error(),
			}
		}
		// Игра не идет.
		if err := a.State(); err != nil {
			return &dto.Status{
				Valid: false,
				Cause: err.Error(),
			}
		}
		if err := a.usecase.Surrender(c.TeamName()); err != nil {
			return &dto.Status{
				Valid: false,
				Cause: err.Error(),
			}
		}
		return &dto.Status{
			Valid: true,
		}
	})
	model.CM.Handle(model.RequestMethodPostNewGame, func(c model.Client, r *dto.Request) any {
		// Вы наблюдатель.
		if c.TeamName() == model.Spectators {
			return &dto.Status{
				Valid: false,
				Cause: model.ErrYouAreSpectator.Error(),
			}
		}
		if err := a.usecase.New(); err != nil {
			return &dto.Status{
				Valid: false,
				Cause: err.Error(),
			}
		}
		return &dto.Status{
			Valid: true,
		}
	})
	model.CM.Handle(model.RequestMethodGetBoard, func(c model.Client, req *dto.Request) any {
		return a.usecase.Export("board")
	})
	model.CM.Handle(model.RequestMethodGetTurn, func(c model.Client, req *dto.Request) any {
		return a.usecase.Export("turn")
	})
	model.CM.Handle(model.RequestMethodGetClock, func(c model.Client, req *dto.Request) any {
		return a.usecase.Export("clock")
	})
	model.CM.Handle(model.RequestMethodGetTeamName, func(c model.Client, req *dto.Request) any {
		return &dto.TeamName{
			Name: c.TeamName(),
		}
	})
	model.CM.Handle(model.RequestMethodGetStatus, func(c model.Client, req *dto.Request) any {
		if a.State() != nil {
			return &dto.Status{
				Valid: false,
				Cause: a.State().Error(),
			}
		}
		return &dto.Status{
			Valid: true,
		}
	})
	model.CM.Handle(model.RequestMethodPostMove, func(c model.Client, req *dto.Request) any {
		var m *dto.Move
		if err := json.Unmarshal(req.Body, &m); err != nil {
			log.Fatal(err)
		}

		// status validation
		if err := a.State(); err != nil {
			return &dto.Status{
				Valid: false,
				Cause: err.Error(),
			}
		}

		// Вы наблюдатель.
		if c.TeamName() == model.Spectators {
			return &dto.Status{
				Valid: false,
				Cause: model.ErrYouAreSpectator.Error(),
			}
		}

		// Не ваш ход.
		if c.TeamName() != a.usecase.Turn() {
			return &dto.Status{
				Valid: false,
				Cause: model.ErrNotYourMove.Error(),
			}
		}

		// usecase validation
		if err := a.usecase.Move(
			entities.NewPosition(m.From.X, m.From.Y),
			entities.NewPosition(m.To.X, m.To.Y),
		); err != nil {
			return &dto.Status{
				Valid: false,
				Cause: err.Error(),
			}
		}

		if err := a.usecase.Mate(); err != nil {
			a.Event("status", &dto.Status{
				Valid: false,
				Cause: err.Error(),
			})
		}

		// Ход валидный и будет сделан.
		return &dto.Status{
			Valid: true,
		}
	})
}
