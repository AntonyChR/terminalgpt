package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/AntonyChR/terminalGPT/color"
)

type Credentials struct {
	Apikey string `json:"api_key"`
}

var userHomeDir, _ = os.UserHomeDir()

var CONFIG_PATH_DIR string = userHomeDir + "/.config/gpt"
var CONFIG_FILE_NAME string = "/gpt.dat"

var CONFIG_FILE_PATH string = path.Join(CONFIG_PATH_DIR, CONFIG_FILE_NAME)

func (c *Credentials) Save() {
	if _, err := os.Stat(CONFIG_PATH_DIR + CONFIG_FILE_NAME); err == nil {
		fmt.Println(color.Yellow("[i] A configuration file already exists, it will be overwritten with the new credentials"))
	}
	os.Mkdir(CONFIG_PATH_DIR, os.ModePerm)
	file, _ := os.Create(CONFIG_PATH_DIR + CONFIG_FILE_NAME)
	jsonBytes, _ := json.Marshal(c)
	file.WriteString(string(jsonBytes))
	fmt.Println(color.Green("[] The credentials have been saved successfully"))
}

func (c *Credentials) Get() error {
	if _, err := os.Stat(CONFIG_PATH_DIR + CONFIG_FILE_NAME); err != nil {
		err := errors.New(color.Red("[] the configuration file does not exist,use the -c flag to add the credentials"))
		return err
	}
	file, _ := os.ReadFile(CONFIG_PATH_DIR + CONFIG_FILE_NAME)
	err := json.Unmarshal(file, &c)

	return err
}

func (c *Credentials) PromptUserForCredentials() error {
	var apiKey string
	fmt.Print("Api key: ")
	fmt.Scanln(&apiKey)
	if apiKey== "" {
		return errors.New("Invalid api key")
	}
	c.Apikey = apiKey
	return nil
}

func (c *Credentials) Exist() bool {
	if c.Apikey != "" {
		return true
	}
	return false
}
