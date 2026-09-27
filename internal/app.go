package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/skvdmt/chess-game-back/internal/delivery"
	"github.com/skvdmt/chess-game-back/internal/model"
)

const (
	defaultTimeout        = 10
	defaultMaxHeaderBytes = 1 << 20 // 1Mb
)

// App Основная структура приложения.
type App struct {
	// Канал сигналов операционной системы для отслеживания
	// сигналов прерывания работы приложения.
	interrupt chan os.Signal
	// Контекст приложения.
	ctx context.Context
	// Функция отмены контекста всего приложения.
	cancel context.CancelFunc
	// Транспортный слой.
	delivery Delivery
	// Роутер.
	router *http.ServeMux
	// Сервер.
	server *http.Server
	// Приложение запущено.
	started bool
	// Ошибки приложения.
	eg []error
	// Приложение уже закрывается.
	stopping bool
	// Корректное завершение горутин.
	wg *sync.WaitGroup
}

// NewApp Конструктор.
func NewApp() (*App, error) {
	model.Logs.Info.Info(fmt.Sprintf("%s creating", model.APP_NAME))
	// Загрузка конфигурации.
	if err := model.CreateConfig(); err != nil {
		return nil, err
	}
	// Создаем глобальный канал ошибок.
	model.Errors = make(chan error)
	// Создание сервера.
	r := http.NewServeMux()
	a := &App{
		interrupt: make(chan os.Signal),
		router:    r,
		server: &http.Server{
			Addr:           fmt.Sprintf(":%d", model.Config.Server.Port),
			Handler:        r,
			ReadTimeout:    defaultTimeout * time.Second,
			WriteTimeout:   defaultTimeout * time.Second,
			MaxHeaderBytes: defaultMaxHeaderBytes,
		},
		wg: &sync.WaitGroup{},
	}

	// Создание контекста.
	a.ctx, a.cancel = context.WithCancel(context.Background())
	// Создание транспортного слоя из которого по
	// цепочки создаются остальные слои приложения.
	var err error
	a.delivery, err = delivery.NewApp(a.ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Start Запуск приложения.
func (a *App) Start() error {
	go a.errorHandler()
	model.Logs.Info.Info(fmt.Sprintf("%s starting", model.APP_NAME))
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		// Настройка и запуск сервера.
		a.routes()
		model.Logs.Info.Info(fmt.Sprintf("http server starting on %d port",
			model.Config.Server.Port))
		if err := a.server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			model.Errors <- err
		}
	}()

	a.wg.Add(1)
	go func() {
		// Запуск слоев приложения по цепочке.
		defer a.wg.Done()
		if err := a.delivery.Start(a.ctx); err != nil {
			model.Errors <- err
		}
		a.started = true
	}()
	return a.interruptHandler()
}

// errorHandle Обработчик глобального канала ошибок.
func (a *App) errorHandler() {
	model.Logs.Info.Info("error handler starting")
	for {
		err := <-model.Errors
		if err == nil {
			return
		}
		if errors.Is(err, context.Canceled) {
			continue
		}
		a.eg = append(a.eg, err)
		if !a.stopping {
			a.interrupt <- syscall.SIGTERM
		}
	}
}

// interruptHandler Обработчик сигналов остановки приложения.
func (a *App) interruptHandler() error {
	model.Logs.Info.Info("error handler starting")
	signal.Notify(a.interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-a.interrupt
	// close sources
	if err := a.stop(); err != nil {
		model.Errors <- err
	}
	model.Errors <- nil
	close(model.Errors)
	if len(a.eg) > 0 {
		return fmt.Errorf("%v", a.eg)
	}
	return nil
}

// Stop Остановка приложения.
func (a *App) stop() error {
	a.stopping = true

	// Отмена контекста.
	a.cancel()
	model.Logs.Info.Info("context canceled")

	// Дождаться запуска приложения.
	for {
		if a.started {
			break
		}
	}

	// Остановка сервера.
	if err := a.server.Shutdown(context.Background()); err != nil {
		return err
	}
	model.Logs.Info.Info("http server shutdown")

	// Остановка транспортного слоя из которого по цепочке
	// останавливаются все остальные слои.
	if err := a.delivery.Stop(a.ctx); err != nil {
		return err
	}
	a.wg.Wait()
	// Закрытие канала отслеживающего сигналы
	// прерывания операционной системы.
	close(a.interrupt)
	model.Logs.Info.Info(fmt.Sprintf("%s stopped", model.APP_NAME))
	return nil
}

// routing Настройка маршрутов.
func (a *App) routes() {
	model.Logs.Info.Info("setup routing")
	a.router.HandleFunc(model.Config.Server.ConnectUrl, a.delivery.Connect)
}
