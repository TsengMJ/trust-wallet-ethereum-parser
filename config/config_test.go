package config_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"ethereum-parser/config"
)

func TestInitConfig(t *testing.T) {
	cases := []struct {
		name       string
		folderPath string
		env        string

		expected    *config.EnvConfig
		expectedErr error
	}{
		{
			name:       "Valid configuration",
			folderPath: "./testdata",
			env:        "test",
			expected: &config.EnvConfig{
				Server: config.Server{
					Port: 8080,
				},
				Log: config.Log{
					Level:             "debug",
					Path:              "./logs",
					Name:              "app.log",
					ErrorLog:          "error.log",
					WarnLog:           "warn.log",
					InfoLog:           "info.log",
					DebugLog:          "debug.log",
					MaxSize:           1,
					MaxAge:            1,
					MaxBackups:        1,
					DisableStacktrace: false,
				},
				Ethereum: config.Ethereum{
					Url: "ethereum-rpc-url",
				},
				Cron: config.Cron{
					Url:    "ethereum-rpc-url",
					Period: "@every 1s",
				},
			},
			expectedErr: nil,
		},
		{
			name:        "Empty folder path",
			folderPath:  "",
			env:         "test",
			expected:    nil,
			expectedErr: errors.New("Folder path is empty"),
		},
		{
			name:        "Non-existent configuration file",
			folderPath:  "./testdata",
			env:         "nonexistent",
			expected:    nil,
			expectedErr: errors.New("Config file does not exist"),
		},
		{
			name:        "Invalid configuration file",
			folderPath:  "./testdata",
			env:         "invalid",
			expected:    nil,
			expectedErr: errors.New("Failed to decode config file"),
		},
		{
			name:        "Empty environment",
			folderPath:  "./testdata",
			env:         "",
			expected:    nil,
			expectedErr: errors.New("Environment is empty"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			config.Config = config.EnvConfig{}
			err := config.InitConfig(&c.folderPath, &c.env)

			if c.name == "Invalid configuration file" {
				fmt.Println(err)
				fmt.Println(config.GetConfig())
			}

			if c.expectedErr != nil {
				assert.Error(t, err, "Expected an error")
				assert.ErrorContains(t, err, c.expectedErr.Error(), "Unexpected error")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.Equal(t, c.expected, config.GetConfig(), "Config does not match expected")
			}
		})
	}
}
