package logger_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"ethereum-parser/config"
	"ethereum-parser/logger"
)

func TestInitLog(t *testing.T) {
	cases := []struct {
		name          string
		logConfig     *config.Log
		expectedErr   error
		expectedPanic bool
	}{
		{
			name: "Valid configuration",
			logConfig: &config.Log{
				Path:     "./logs",
				ErrorLog: "error.log",
				WarnLog:  "warn.log",
				InfoLog:  "info.log",
				DebugLog: "debug.log",
				MaxAge:   7, // 7 days
			},
			expectedErr:   nil,
			expectedPanic: false,
		},
		{
			name: "Empty log path",
			logConfig: &config.Log{
				Path:     "",
				ErrorLog: "error.log",
				WarnLog:  "warn.log",
				InfoLog:  "info.log",
				DebugLog: "debug.log",
				MaxAge:   7, // 7 days
			},
			expectedErr:   errors.New("Log path is empty"),
			expectedPanic: false,
		},
		{
			name: "Empty error log",
			logConfig: &config.Log{
				Path:     "./logs",
				ErrorLog: "",
				WarnLog:  "warn.log",
				InfoLog:  "info.log",
				DebugLog: "debug.log",
				MaxAge:   7, // 7 days
			},
			expectedErr:   errors.New("Error log is empty"),
			expectedPanic: false,
		},
		{
			name: "Empty warn log",
			logConfig: &config.Log{
				Path:     "./logs",
				ErrorLog: "error.log",
				WarnLog:  "",
				InfoLog:  "info.log",
				DebugLog: "debug.log",
				MaxAge:   7, // 7 days
			},
			expectedErr:   errors.New("Warn log is empty"),
			expectedPanic: false,
		},
		{
			name: "Empty info log",
			logConfig: &config.Log{
				Path:     "./logs",
				ErrorLog: "error.log",
				WarnLog:  "warn.log",
				InfoLog:  "",
				DebugLog: "debug.log",
				MaxAge:   7, // 7 days
			},
			expectedErr:   errors.New("Info log is empty"),
			expectedPanic: false,
		},
		{
			name: "Empty debug log",
			logConfig: &config.Log{
				Path:     "./logs",
				ErrorLog: "error.log",
				WarnLog:  "warn.log",
				InfoLog:  "info.log",
				DebugLog: "",
				MaxAge:   7, // 7 days
			},
			expectedErr:   errors.New("Debug log is empty"),
			expectedPanic: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Set the configuration
			config.Config.Log = *c.logConfig

			// Cleanup logs directory after test
			defer os.RemoveAll("./logs")

			// Call InitLog and capture panic if any
			defer func() {
				if r := recover(); r != nil {
					assert.True(t, c.expectedPanic, "Unexpected panic occurred")
				}
			}()

			err := logger.InitLog()

			if c.expectedErr != nil {
				assert.Error(t, err, "Expected an error")
				assert.ErrorContains(t, err, c.expectedErr.Error(), "Unexpected error")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.NotNil(t, logger.Logger, "Logger is nil after initialization")
			}
		})
	}
}
