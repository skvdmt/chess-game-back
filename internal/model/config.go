package model

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	// Название приложения.
	APP_NAME = "chess-back-game"

	// Путь в директории конфигурации. (Добавляется директория с именем приложения).
	configDirectoryProd = "/etc"
	configDirectoryDev  = "./config"
	// Имя файла конфигурации.
	configFileNameProd = "prod.yaml"
	configFileNameDev  = "dev.yaml"
)

// Config Глобальная конфигурация.
var Config *MainConfig

// ServerConfig Конфигурация сервера.
type ServerConfig struct {
	Port              uint16 `yaml:"port"`
	ConnectUrl        string `yaml:"connect_url"`
	OriginalClientUrl string `yaml:"original_client_url"`
}

type GameConfig struct {
	StepTimeLeft                int  `yaml:"step_time_left"`
	ReserveTimeLeft             int  `yaml:"reserve_time_left"`
	OfferDrawTimesLeft          int  `yaml:"offer_draw_times_left"`
	TimeLeftForConfirmDraw      int  `yaml:"time_left_for_confirm_draw"`
	SwapTeamsAfterMakingNewGame bool `yaml:"swap_teams_after_making_new_game"`
}

// MainConfig Основная конфигурация.
type MainConfig struct {
	Server *ServerConfig `yaml:"server"`
	Game   *GameConfig   `yaml:"game"`
}

// CreateConfig Загрузка конфигурации в глобальную переменную Config.
func CreateConfig() error {
	Logs.Info.Info("configuration loading")
	dn := filepath.Join(configDirectoryProd, APP_NAME)
	fn := configFileNameProd
	m, ok := os.LookupEnv(MODE)
	if ok && m == Dev {
		dn = configDirectoryDev
		fn = configFileNameDev
	}
	d, err := os.ReadFile(filepath.Join(dn, fn))
	if err != nil {
		return err
	}
	cfg := &MainConfig{}
	if err := yaml.Unmarshal(d, cfg); err != nil {
		return err
	}
	Config = cfg
	return nil
}
