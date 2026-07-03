package model

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	MODE = "MODE"
	Dev  = "dev"
	// Путь к директории журналов. (Добавляется директория с именем приложения).
	logDirectoryProd = "/var/log"
	logDirectoryDev  = "./logs"
	// Имя файла журнала ошибок.
	logFileName = "error.log"
	logFlag     = os.O_CREATE | os.O_APPEND | os.O_RDWR
	logPerm     = 0666
)

// Logs Глобальная переменная инструмента медения журнала.
var Logs *Logger

// logger Инструмент ведения журнала.
type Logger struct {
	// Журнал информирования.
	Info *slog.Logger
	// Журнал ошибок.
	Error *slog.Logger
	// Файл журнала ошибок.
	errorFile *os.File
}

// Close Закрытие ресурсов логгера.
func (l *Logger) Close() error {
	if l.errorFile != nil {
		return l.errorFile.Close()
	}
	return nil
}

// Createlogger Создать логгер и установить ссылку на него
// в глобальную переменную Logs. В логгере создается зеркало
// ошибок в os.Stderr и файл журнала.
func CreateLogger() error {
	// Установка директории файлов журнала.
	logDirectory := logDirectoryProd
	mode, ok := os.LookupEnv(MODE)
	if ok && mode == Dev {
		logDirectory = logDirectoryDev
	}
	n := "models.logger.Loadlogger"
	// Создать дерикторию журнала для приложения в случае ее отсутствия.
	dn := filepath.Join(logDirectory, APP_NAME)
	if _, err := os.Stat(dn); os.IsNotExist(err) {
		if err := os.MkdirAll(dn, os.ModePerm); err != nil {
			return err
		}
	}
	// Открыть файл журнала ошибок.
	fn := filepath.Join(logDirectory, APP_NAME, logFileName)
	ef, err := os.OpenFile(fn, logFlag, logPerm)
	if err != nil {
		return fmt.Errorf("%s %w", n, err)
	}
	// Писать журнал ошибок в os.Stderr, а также в файл.
	ew := io.MultiWriter(os.Stderr, ef)
	// Создать логгер в глобальной переменной Logs.
	Logs = &Logger{
		Info:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
		Error:     slog.New(slog.NewJSONHandler(ew, nil)),
		errorFile: ef,
	}
	return nil
}
