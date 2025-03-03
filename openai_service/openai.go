package openaiservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	color "github.com/AntonyChR/terminalGPT/colors"
	"github.com/AntonyChR/terminalGPT/logs"
)

var ChatRoles = Roles{
	System:    "system",
	Assistant: "assistant",
	User:      "user",
}

// prompt to generate title, file name with title plus txt extension, this will be generated in json format
var PROMPT_TO_GET_TITLE_EN string = `Generate a title for the current conversation and a corresponding file name with a .md extension. 
The file name must be the same as the title, replacing spaces with underscores. The title should be descriptive and concise, written in English. 
The response should be generated in JSON format in one line.

Example:
Title: "My first conversation"
File name: "My_first_conversation.md"

the response looks like this:
{"title": "My first conversation","file_name": "My_first_conversation.md"}
`

var PROMPT_TO_GET_TITLE_ES string = `Genera un título para la conversación actual y un nombre de archivo correspondiente con una extensión .md.
El nombre del archivo debe ser igual al título, reemplazando los espacios con guiones bajos. El título debe ser descriptivo y conciso.
La respuesta debe generarse en formato JSON en una sola línea.

Ejemplo:
Título: "Mi primera conversación"
Nombre de archivo: "Mi_primera_conversación.md"

la repsuesta que generas debería verse así:
{"title": "Mi primera conversación","file_name": "Mi_primera_conversación.md"}
`

type ApiConfiguration struct {
	Apikey        string
	Model         string
	ApiUrl        string
	InitialPrompt string
	Temperature   float32 //default 0.0
	Stream        bool    //default false
}

func NewChat(c ApiConfiguration) *Openai {
	initialMessage := Message{
		Role:    ChatRoles.System,
		Content: c.InitialPrompt,
	}
	return &Openai{
		apikey: c.Apikey,
		url:    c.ApiUrl,
		Chat: Chat{
			Model: c.Model,
			Messages: []Message{
				initialMessage,
			},
			Temperature: c.Temperature,
			Stream:      c.Stream,
		},
	}
}

type Openai struct {
	apikey     string
	url        string
	Chat       Chat
	LogChannel chan string
}

type Chat struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Temperature float32   `json:"temperature,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

func (c *Chat) appendMessage(m Message) {
	c.Messages = append(c.Messages, m)

}

func (o *Openai) AddMessageAsUser(input string) {
	o.Chat.appendMessage(Message{Role: ChatRoles.User, Content: input})
}

func (o *Openai) AddMessageAsAssistant(input string) {
	o.Chat.appendMessage(Message{Role: ChatRoles.Assistant, Content: input})
}

func (o *Openai) CreateRequest() (*http.Request, error) {
	encodedData, _ := json.Marshal(o.Chat)
	req, err := http.NewRequest("POST", o.url, bytes.NewBuffer(encodedData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+o.apikey)
	req.Header.Set("Content-Type", "application/json")
	return req, nil

}

func (o *Openai) RemoveLastMessage() {
	o.Chat.Messages = o.Chat.Messages[:len(o.Chat.Messages)-1]
}

func (o *Openai) GetCompletion() (string, error) {
	o.Chat.Stream = false
	req, err := o.CreateRequest()
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		o.LogChannel <- createRequestErrorEvent(resp, err, "")
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 429 {
			o.LogChannel <- createRequestErrorEvent(resp, err, "Request limit exceeded")
			panic("Request limit exceeded")
		}
		o.LogChannel <- createRequestErrorEvent(resp, err, "")
		panic("unespected request error")
	}

	var completionResponse CompletionResponse
	bodyBytes, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &completionResponse)
	if err != nil {
		return "", fmt.Errorf("error Unmarshal: %v", err)
	}

	return completionResponse.Choices[0].Message.Content, nil

}

// get completion with stream
func (o *Openai) GetStreamCompletion(stringChannel chan string) error {
	o.Chat.Stream = true
	req, _ := o.CreateRequest()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("error making request " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 429 {
			panic("Request limit exceeded")
		}
		panic("unespected request error")
	}

	chunkBuffer := make([]byte, 4096)
	re := regexp.MustCompile(`(?s)"finish_reason":(null|"stop").*`)
	msg := Message{Role: ChatRoles.Assistant}
	content := ""
	for {
		chunkSize, err := resp.Body.Read(chunkBuffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				stringChannel <- "\n\n"
				break
			}
			logStr := "error reading stream completion chunk " + err.Error()
			fmt.Println(logStr)
			o.LogChannel <- logs.NewEvent(logs.ErrorLevel, logStr, logs.GetDate()).ToStr()
		}

		chunkString := string(chunkBuffer[:chunkSize])
		sep := `data: {"id":"`

		splitedChunkString := strings.Split(chunkString, sep)

		for index, str := range splitedChunkString {
			if index == 0 {
				continue
			}
			var chunkObject CompletionResponse
			validJson := `{"id":"` + re.ReplaceAllString(str, `"finish_reason":"stop"}]}`)
			err := json.Unmarshal([]byte(validJson), &chunkObject)
			if err != nil {
				logStr := "error unmarshalling json del stream completion chunk " + err.Error()
				fmt.Print(color.Red(logStr))
				o.LogChannel <- logs.NewEvent(logs.ErrorLevel, logStr, logs.GetDate()).ToStr()
				continue
			}
			deltaString := chunkObject.Choices[0].Delta.Content
			content += deltaString
			stringChannel <- color.Green(deltaString)
		}

	}

	msg.Content = content
	o.Chat.appendMessage(msg)
	return nil

}

func (o *Openai) Reset() {
	fmt.Println(color.Red("Reset context"))
	o.Chat.Messages = []Message{o.Chat.Messages[0]}
}

// save chat as md file
func (o *Openai) Save(title, fileName, dir string) error {

	//check if the route exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	route := path.Join(dir, fileName)

	file, err := os.Create(route)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("# Title: " + title + "\n\n")

	for _, message := range o.Chat.Messages {
		if message.Role == ChatRoles.User {
			file.WriteString("### [A] " + message.Content + "\n\n---\n\n")
		}
		if message.Role == ChatRoles.Assistant {
			file.WriteString(message.Content)
		}
	}
	return nil

}

func createRequestErrorEvent(resp *http.Response, err error, additionalInfo string) string {
	bodyBytes, _ := io.ReadAll(resp.Body)
	logMsg := fmt.Sprintf("error making request\n%s \n%s\nResponse body: %s\nStatus code: %d", additionalInfo, err.Error(), string(bodyBytes), resp.StatusCode)
	event := logs.NewEvent(logs.ErrorLevel, logMsg, logs.GetDate())
	return event.ToStr()
}
