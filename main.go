package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/AntonyChR/terminalGPT/color"
	"github.com/AntonyChR/terminalGPT/config"
	openaiservice "github.com/AntonyChR/terminalGPT/openai_service"
)

func readUserInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("[A]: ")
	input, _ := reader.ReadString('\n')
	return input, validateInput(input)
}

func validateInput(input string) error {
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
		Apikey: credentials.Apikey,
		Model:  "gpt-3.5-turbo",
		ApiUrl: "https://api.openai.com/v1/chat/completions",
	}
	chat := openaiservice.NewChat(config)

	args := os.Args
	var input string
	if len(args) == 2 {
		input = args[1]
	}

	if input != "" {
		chat.AddMessageAsUser(input)
		completion, err := chat.GetCompletion()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(completion.Content)
		os.Exit(0)
	}

	var err error

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
		fmt.Println(color.Green("[GPT]: "+completion.Content) + "\n")
	}

}
