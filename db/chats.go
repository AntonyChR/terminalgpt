package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	openaiservice "github.com/AntonyChR/terminalGPT/openai_service"
)

var DB_FILE_NAME = "/chat.json"

var USER_HOME_DIR, _ = os.UserHomeDir()
var DB_FOLDER = USER_HOME_DIR + "/gpt"
var DB_FILE = path.Join(DB_FOLDER, DB_FILE_NAME)

func SaveChat(title string, chat openaiservice.Chat) error {
	if _, err := os.Stat(DB_FOLDER); err != nil {
		err = os.Mkdir(DB_FOLDER, 0777)
		if err != nil {
			return err
		}
	}

	dataBytes, _ := json.Marshal(chat)

	err := ioutil.WriteFile(DB_FILE, dataBytes, 0644)

	return err

}
