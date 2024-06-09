package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type EnvConfig struct {
	Server   Server   `toml:"server"`
	Log      Log      `toml:"log"`
	Ethereum Ethereum `toml:"ethereum"`
	Cron     Cron     `toml:"cron"`
}

type Server struct {
	Port int `toml:"port"`
}

type Log struct {
	Level             string `toml:"level"`
	Path              string `toml:"path"`
	Name              string `toml:"name"`
	ErrorLog          string `toml:"error_log"`
	WarnLog           string `toml:"warn_log"`
	InfoLog           string `toml:"info_log"`
	DebugLog          string `toml:"debug_log"`
	MaxSize           int    `toml:"max_size"`
	MaxAge            int    `toml:"max_age"`
	MaxBackups        int    `toml:"max_backups"`
	DisableStacktrace bool   `toml:"disable_stacktrace"`
}

type Ethereum struct {
	Url string `toml:"url"`
}

type Cron struct {
	Url    string `toml:"url"`
	Period string `toml:"period"`
}

var Config EnvConfig

func InitConfig(folderPath *string, env *string) error {
	if *folderPath == "" {
		return errors.New("Folder path is empty")
	}

	if *env == "" {
		return errors.New("Environment is empty")
	}

	filePath := fmt.Sprintf("%s/%s.toml", *folderPath, *env)
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return errors.New("Config file does not exist, " + err.Error())
	}

	_, err := toml.DecodeFile(filePath, &Config)
	if err != nil {
		return errors.New("Error decoding config file, " + err.Error())
	}

	return nil
}

func GetConfig() *EnvConfig {
	return &Config
}
