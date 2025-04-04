package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"

	colors "github.com/AntonyChR/terminalGPT/colors"
	configService "github.com/AntonyChR/terminalGPT/config"
	openaiservice "github.com/AntonyChR/terminalGPT/openai_service"
	userinterface "github.com/AntonyChR/terminalGPT/user_interface"
	utils "github.com/AntonyChR/terminalGPT/utils"
)

const (
	CONFIG_FILE_JSON string = "/config.json"
	LOG_FILE         string = "/logs.txt"
)

var (
	// paths
	CONFIG_DIR        string = path.Join(HOME, ".config/gpt")
	CONVERSATIONS_DIR string = path.Join(HOME, "/conversations")
	HOME, _                  = os.UserHomeDir()
	CURRENT_DATE      string = utils.GetCurrentDate()

	// flags
	lang = flag.String("lang", "es", "-lang <es,en>")
	file = flag.String("file", "", "-file <file_name>")
)

type FileTitle struct {
	Title    string `json:"title"`
	FileName string `json:"file_name"`
}

func main() {
	flag.Parse()

	userInterface := userinterface.New()
	go userInterface.PrintChannelMessages()

	configFile := path.Join(CONFIG_DIR, CONFIG_FILE_JSON)

	config, err := configService.LoadConfig(configFile)

	if err != nil {
		handleError(err, "Error loading config file")
		config, err = createNewConfig(userInterface)
		handleError(err, "Error creating new config file")

		err = configService.SaveNewConfig(config, CONFIG_DIR, "config.json")
		handleError(err, "Error saving new config file")
	}

	if *lang != "" && *lang != config.Lang {
		println("Changing language to " + *lang)
		config.Lang = *lang
		err = configService.SaveNewConfig(config, CONFIG_DIR, "config.json")
		handleError(err, "Error saving new config file")
	}

	fileContent := ""
	if *file != "" {
		//read content as string
		content, err := os.ReadFile(*file)
		handleError(err, "Error reading file")
		fileContent = string(content)

		fmt.Printf("\nThe content of the file will be concatenated to the prompt\n")
		fmt.Printf("\n%s:\n\n%s\n...\n\n",*file,fileContent[:200])
		//print first lines of the file 

	}

	sessionConfig := openaiservice.ApiConfiguration{
		Apikey:        config.ApiKey,
		ApiUrl:        "https://api.deepseek.com/chat/completions",
		Model:         "deepseek-chat",
		InitialPrompt: config.InitialPrompt,
	}

	chat := openaiservice.NewChat(sessionConfig)

	for {
		input, err := userInterface.GetInput(colors.Yellow("[A]: "))

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		switch input {

		case "reset":
			chat.Reset()
			continue
		case "save":
			if config.Lang == configService.SPANISH {
				chat.AddMessageAsUser(openaiservice.PROMPT_TO_GET_TITLE_ES)
			} else {
				chat.AddMessageAsUser(openaiservice.PROMPT_TO_GET_TITLE_EN)
			}
			println("Generating title...")
			comp, err := chat.GetCompletion()
			handleError(err, "Error getting completion")

			// remove completion prompt
			chat.RemoveLastMessage()

			comp = utils.FixJsonString(comp)

			var fileTitle FileTitle

			err = json.Unmarshal([]byte(comp), &fileTitle)
			handleError(err, "Error unmarshalling json")

			err = chat.Save(fileTitle.Title, fileTitle.FileName, CONVERSATIONS_DIR)
			handleError(err, "Error saving conversation")

			fmt.Println(colors.Green("Conversation saved successfully in " + CONVERSATIONS_DIR + "/" + fileTitle.FileName))

			// remove info generated
			chat.RemoveLastMessage()
			continue

		case "exit":
			os.Exit(0)
		}

		chat.AddMessageAsUser(input + "\n" + fileContent)
		err = chat.GetStreamCompletion(userInterface.PrintChannel)
		handleError(err, "Error getting completion")
	}

}

func createNewConfig(userInt *userinterface.UserInterfaceIO) (*configService.Config, error) {
	userInt.PrintChannel <- colors.Yellow("Please enter your API key: ")
	api_key, err := userInt.GetInput(colors.Yellow(">: "))

	handleError(err, "Error getting input")

	newConf := &configService.Config{
		ApiKey: api_key,
	}

	lang, err := userInt.GetInput(colors.Yellow("Language (es/en): "))
	handleError(err, "Error getting input")

	if lang != configService.SPANISH && lang != configService.ENGLISH {
		println(colors.Red("Invalid language, spanish will be set by default"))
		lang = configService.SPANISH
	}

	newConf.Lang = lang

	newConf.InitialPrompt = "You are an helpful assistant"

	newConf.InitialPrompt = newConf.InitialPrompt + ", the current date is " + CURRENT_DATE
	if err != nil {
		return nil, err
	}

	return newConf, nil
}

func handleError(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %v\n", message, err)
		os.Exit(1)
	}
}
