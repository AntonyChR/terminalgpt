package openaiservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	color "github.com/AntonyChR/terminalGPT/colors"
	customErrors "github.com/AntonyChR/terminalGPT/customErrors"
)

var ChatRoles = Roles{
	System:    "system",
	Assistant: "assistant",
	User:      "user",
}

type ApiConfiguration struct {
	Apikey      string
	Model       string
	ApiUrl      string
	Temperature float32 //default 0.0
}

func NewChat(c ApiConfiguration) *Openai {
	initialMessage := Message{
		Role:    ChatRoles.System,
		Content: "You are an assistant",
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
			Stream:      true,
		},
	}
}

type Openai struct {
	apikey string
	url    string
	Chat   Chat
}

type Chat struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Temperature float32   `json:"temperature,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

func (o *Openai) appendMessageToList(m Message) {
	o.Chat.Messages = append(o.Chat.Messages, m)
}

func (o *Openai) AddMessageAsUser(input string) {
	o.appendMessageToList(Message{Role: ChatRoles.User, Content: input})
}

func (o *Openai) GetStreamCompletion(streamChannel chan string) error {
	encodedData, _ := json.Marshal(o.Chat)
	req, err := http.NewRequest("POST", o.url, bytes.NewBuffer(encodedData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+o.apikey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 429 {
			return &customErrors.RequestError{StatusCode: resp.StatusCode, Err: errors.New("Request limit exceeded")}
		}
		return &customErrors.RequestError{StatusCode: resp.StatusCode, Err: errors.New("unespected request error")}
	}

	chunkBuffer := make([]byte, 4096)
	re := regexp.MustCompile(`(?s)"finish_reason":(null|"stop").*`)
	msg := Message{Role: ChatRoles.Assistant}
	content := ""
	fmt.Printf("\n")
	for {
		chunkSize, err := resp.Body.Read(chunkBuffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				streamChannel <- "\n\n"
				break
			}
			fmt.Println("ERROR: ", err.Error())
		}

		chunkString := string(chunkBuffer[:chunkSize])

		sep := `data: {"id":"chatcmpl`

		splitedChunkString := strings.Split(chunkString, sep)

		for index, str := range splitedChunkString {
			if index == 0 {
				continue
			}
			var chunkObject CompletionResponse
			validJson := `{"id":"chatcmpl` + re.ReplaceAllString(str, `"finish_reason":"stop"}]}`)
			err := json.Unmarshal([]byte(validJson), &chunkObject)
			if err != nil {
				fmt.Println(color.Colorize("red", err.Error()))
				continue
			}
			deltaString := chunkObject.Choices[0].Delta.Content
			content += deltaString
			streamChannel <- color.Colorize("green", deltaString) 
		}

	}

	msg.Content = content
	o.appendMessageToList(msg)
	return nil

}

func (o *Openai) Reset() {
	fmt.Println(color.Colorize("red", "Reset context"))
	o.Chat.Messages = []Message{o.Chat.Messages[0]}
}
