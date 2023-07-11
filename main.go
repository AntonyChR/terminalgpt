package main

import (
	"fmt"
	"os"

	config "github.com/AntonyChR/terminalGPT/config"
	openaiservice "github.com/AntonyChR/terminalGPT/openai_service"
	userinterface "github.com/AntonyChR/terminalGPT/user_interface"
)

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

	config := openaiservice.ApiConfiguration{
		Apikey: credentials.Apikey,
		Model:  "gpt-3.5-turbo",
		ApiUrl: "https://api.openai.com/v1/chat/completions",
	}
	userInterface := userinterface.NewUserInterface()
	chat := openaiservice.NewChat(config)

	var input string
	var err error

	go userInterface.PrintChannelMessages()

	for {
		input, err = userInterface.GetInput("[A]: ")
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if input == "reset\n" {
			chat.Reset()
			continue
		}
		chat.AddMessageAsUser(input)
		completion, err := chat.GetStreamCompletion(userInterface.PrintChannel)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		chat.AddMessage(completion)
	}

}
