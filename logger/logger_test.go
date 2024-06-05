package logger_test

import (
	"ethereum-parser/config"
	"ethereum-parser/logger"
	"testing"
)

func TestInitLog(t *testing.T) {
	// Test Case 0: Success
	t.Run("Success", func(t *testing.T) {
		config.Config = config.EnvConfig{
			Log: config.Log{
				Path:     "./test_logs",
				ErrorLog: "error.log",
				WarnLog:  "warn.log",
				InfoLog:  "info.log",
				DebugLog: "debug.log",
				MaxAge:   7,
			},
		}
		err := logger.InitLog()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	// Test Case 1: Log path is empty
	t.Run("Log path is empty", func(t *testing.T) {
		config.Config = config.EnvConfig{
			Log: config.Log{
				Path:     "",
				ErrorLog: "error.log",
				WarnLog:  "warn.log",
				InfoLog:  "info.log",
				DebugLog: "debug.log",
				MaxAge:   7,
			},
		}
		err := logger.InitLog()
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	// Test Case 2: Error log is empty
	t.Run("Error creating log directory", func(t *testing.T) {
		config.Config = config.EnvConfig{
			Log: config.Log{
				Path:     "./test_logs",
				ErrorLog: "",
				WarnLog:  "warn.log",
				InfoLog:  "info.log",
				DebugLog: "debug.log",
				MaxAge:   7,
			},
		}
		err := logger.InitLog()
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	// Test Case 3: Warn log is empty
	t.Run("Warn log is empty", func(t *testing.T) {
		config.Config = config.EnvConfig{
			Log: config.Log{
				Path:     "./test_logs",
				ErrorLog: "error.log",
				WarnLog:  "",
				InfoLog:  "info.log",
				DebugLog: "debug.log",
				MaxAge:   7,
			},
		}
		err := logger.InitLog()
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	// Test Case 4: Info log is empty
	t.Run("Info log is empty", func(t *testing.T) {
		config.Config = config.EnvConfig{
			Log: config.Log{
				Path:     "./test_logs",
				ErrorLog: "error.log",
				WarnLog:  "warn.log",
				InfoLog:  "",
				DebugLog: "debug.log",
				MaxAge:   7,
			},
		}
		err := logger.InitLog()
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	// Test Case 5: Debug log is empty
	t.Run("Debug log is empty", func(t *testing.T) {
		config.Config = config.EnvConfig{
			Log: config.Log{
				Path:     "./test_logs",
				ErrorLog: "error.log",
				WarnLog:  "warn.log",
				InfoLog:  "info.log",
				DebugLog: "",
				MaxAge:   7,
			},
		}
		err := logger.InitLog()
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

}
