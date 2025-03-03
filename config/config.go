package config

import (
	"encoding/json"
	"errors"
	"os"
	"path"
)

var SPANISH = "es"
var ENGLISH = "en"

type Config struct {
	ApiKey        string `json:"api_key,omitempty"`
	InitialPrompt string `json:"initial_prompt,omitempty"`
	Lang          string `json:"lang,omitempty"`
}

func SaveNewConfig(config *Config, dir, fileName string) error {

	//check if the route exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}
	route := path.Join(dir, fileName)

	file, err := os.Create(route)
	if err != nil {
		return err
	}
	defer file.Close()

	configJson, err := json.Marshal(config)
	if err != nil {
		return err
	}

	file.Write(configJson)

	return nil

}

func LoadConfig(fileRoute string) (*Config, error) {
	file, err := os.ReadFile(fileRoute)
	config := &Config{}
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, config)
	if err != nil {
		return config, err
	}

	if config.ApiKey == "" {
		return config, errors.New("api_key field is empty")
	}

	return config, nil
}
