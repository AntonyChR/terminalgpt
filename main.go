package main

import (
	"errors"
	"fmt"
	"os"
	config "github.com/AntonyChR/terminalGPT/config"
	openaiservice "github.com/AntonyChR/terminalGPT/openai_service"
)

func readUserInput() (string, error) {
	var input string
	fmt.Scanln(&input)
	return input, validateInput(input)
}

func validateInput(input string) error {
	if input == "" || input == "\n" {
		return errors.New("invalid input")
	}
	return nil
}

func main() {

	var credentials config.Credentials
	if err := credentials.Get(); err != nil {
		err = credentials.PromptUserForCredentials()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		credentials.Save()
	}

	config := openaiservice.Configuration{
		Apikey:     credentials.Apikey,
		Model:      "gpt-3.5-turbo",
		ApiUrl:     "https://api.openai.com/v1/chat/completions",
		StreamData: true,
	}
	chat := openaiservice.NewChat(config)

	var input string
	var err error

	go chat.ListenAndPrintIncommingMsg()
	chat.IncommingMessage <- "\n[A]: "

	for {
		input, err = readUserInput()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if input == "reset\n" {
			chat.Reset()
			continue
		}
		chat.AddMessageAsUser(input)
		completion, err := chat.GetCompletion()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		chat.AddMessage(completion)
	}

}
