package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	cli "github.com/AntonyChR/terminalGPT/cli"
	"github.com/AntonyChR/terminalGPT/color"
	"github.com/AntonyChR/terminalGPT/config"
	openaiservice "github.com/AntonyChR/terminalGPT/openai_service"
)

var (
	INPUT       = flag.String("m", "", cli.M)
	CREDENTIALS = flag.String("c", "", cli.C)
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
	flag.Parse()
	var credentials config.Credentials

	if *CREDENTIALS != "" {
		credentials.Apikey = *CREDENTIALS
		credentials.Save()
	}

	if !credentials.Exist() {
		if err := credentials.Get(); err != nil {
			fmt.Println(err.Error())
			flag.Usage()
			os.Exit(1)
		}
	}

	config := openaiservice.Configuration{
		Apikey: credentials.Apikey,
		Model:  "gpt-3.5-turbo",
		ApiUrl: "https://api.openai.com/v1/chat/completions",
	}
	chat := openaiservice.NewChat(config)

	if *INPUT != "" {
		chat.AddMessageAsUser(*INPUT)
		completion, err := chat.GetCompletion()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(completion.Content)
		os.Exit(0)
	}

	var err error

	for *INPUT != ".exit" {
		*INPUT, err = readUserInput()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		chat.AddMessageAsUser(*INPUT)
		completion, err := chat.GetCompletion()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		chat.AddMessage(completion)
		fmt.Println(color.Green("[GPT]: "+completion.Content) + "\n")
	}

}
