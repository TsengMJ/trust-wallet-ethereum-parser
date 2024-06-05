package config_test

import (
	"ethereum-parser/config"
	"testing"
)

func TestInitConfig(t *testing.T) {
	// Test Case 1: Valid file path and environment
	t.Run("Valid file path and environment", func(t *testing.T) {

		folderPath := "./testdata"
		env := "valid"
		err := config.InitConfig(&folderPath, &env)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	// Test Case 2: Invalid file path
	t.Run("Invalid file path", func(t *testing.T) {
		folderPath := "./invalidpath"
		env := "valid"
		err := config.InitConfig(&folderPath, &env)
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
	})

	// Test Case 3: Invalid environment file
	t.Run("Invalid environment file", func(t *testing.T) {
		folderPath := "./testdata"
		env := "invalid"
		err := config.InitConfig(&folderPath, &env)

		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
	})

	// Test Case 4: Empty folder path
	t.Run("Empty folder path", func(t *testing.T) {
		folderPath := ""
		env := "valid"
		err := config.InitConfig(&folderPath, &env)
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
	})

	// Test Case 5: Empty environment
	t.Run("Empty environment", func(t *testing.T) {
		folderPath := "./testdata"
		env := ""
		err := config.InitConfig(&folderPath, &env)
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
	})

	// Test Case 6: Both folder path and environment empty
	t.Run("Both folder path and environment empty", func(t *testing.T) {
		folderPath := ""
		env := ""
		err := config.InitConfig(&folderPath, &env)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
	})
}
