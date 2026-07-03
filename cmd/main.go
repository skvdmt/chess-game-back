package main

import (
	"os"

	"github.com/skvdmt/chess-game-back/internal"
	"github.com/skvdmt/chess-game-back/internal/model"
)

// main Точка входа в приложение.
func main() {
	// Создание логгера.
	if err := model.CreateLogger(); err != nil {
		panic(err)
	}
	// Создание приложения.
	app, err := internal.NewApp()
	if err != nil {
		model.Logs.Error.Error(err.Error())
		os.Exit(1)
	}
	// Запуск приложения.
	if err := app.Start(); err != nil {
		model.Logs.Error.Error(err.Error())
		os.Exit(1)
	}
	// Закрытие логгера.
	if err := model.Logs.Close(); err != nil {
		panic(err)
	}
}
