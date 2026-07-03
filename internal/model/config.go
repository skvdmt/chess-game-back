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
	configFileNameProd = "config.yaml"
	configFileNameDev  = "config-dev.yaml"
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

// LoadConfig Загрузка конфигурации в глобальную переменную Config.
func LoadConfig() error {
	Logs.Info.Info("configuration loading")
	configDirectory := configDirectoryProd
	configFileName := configFileNameProd
	mode, ok := os.LookupEnv(MODE)
	if ok && mode == Dev {
		configDirectory = configDirectoryDev
		configFileName = configFileNameDev
	}
	d, err := os.ReadFile(filepath.Join(configDirectory, APP_NAME, configFileName))
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
