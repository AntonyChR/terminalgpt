package config

import (
	"os"
	"path"
	"testing"
)

var TEST_CONFIG_PATH_DIR = "./test_config"
var TEST_CONFIG_FILE = "config.json"

func TestCreateAndDeleteConfig(t *testing.T) {
	// Create a new API key and save it
	API_KEY := "test_api_key"
	INITIAL_PROMPT := "test_initial_prompt"

	config := Config{ApiKey: API_KEY, InitialPrompt: INITIAL_PROMPT}
	err := SaveNewConfig(&config, TEST_CONFIG_PATH_DIR, TEST_CONFIG_FILE)
	if err != nil {
		t.Fatalf("Failed to save new API key: %v", err)
	}

	// Verify that the file was created
	route := path.Join(TEST_CONFIG_PATH_DIR, TEST_CONFIG_FILE)
	if _, err := os.Stat(route); os.IsNotExist(err) {
		t.Fatalf("Config file was not created")
	}

	// Load the config from the file
	loadedConfig, err := LoadConfig(route)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify that the loaded config matches the saved config
	if loadedConfig.ApiKey != config.ApiKey {
		t.Errorf("Expected API key %v, got %v", config.ApiKey, loadedConfig.ApiKey)
	}

	// Clean up by deleting the test config file and directory
	err = os.Remove(route)
	if err != nil {
		t.Fatalf("Failed to delete config file: %v", err)
	}

	err = os.Remove(TEST_CONFIG_PATH_DIR)
	if err != nil {
		t.Fatalf("Failed to delete config directory: %v", err)
	}
}
