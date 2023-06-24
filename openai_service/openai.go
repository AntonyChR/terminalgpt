package openaiservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	customErrors "github.com/AntonyChR/terminalGPT/customErrors"
)

var ChatRoles = Roles{
	System:    "system",
	Assistant: "assistant",
	User:      "user",
}

type Configuration struct {
	Apikey      string
	Model       string
	ApiUrl      string
	Temperature float32 //default 0.0
}

func NewChat(c Configuration) *Openai {
	initialMessage := Message{
		Role:    ChatRoles.System,
		Content: "You are an assistant.",
	}
	return &Openai{
		apikey: c.Apikey,
		url:    c.ApiUrl,
		chat: Chat{
			Model: c.Model,
			Messages: []Message{
				initialMessage,
			},
			Temperature: c.Temperature,
		},
	}
}

type Openai struct {
	apikey string
	url    string
	chat   Chat
}

type Chat struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Temperature float32   `json:"temperature,omitempty"`
}

func (o *Openai) AddMessage(m Message) {
	o.chat.Messages = append(o.chat.Messages, m)
}

func (o *Openai) AddMessageAsUser(input string) {
	o.AddMessage(Message{Role: ChatRoles.User, Content: input})
}

func (o *Openai) GetCompletion() (Message, error) {
	encodedData, _ := json.Marshal(o.chat)
	req, err := http.NewRequest("POST", o.url, bytes.NewBuffer(encodedData))
	if err != nil {
		return Message{}, err
	}
	req.Header.Add("Authorization", "Bearer "+o.apikey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Message{}, &customErrors.RequestError{StatusCode: resp.StatusCode, Err: err}
	}

	if resp.StatusCode != http.StatusOK {
		return Message{}, &customErrors.RequestError{StatusCode: resp.StatusCode, Err: errors.New("unespected request errror")}
	}

	completionResp := CompletionResponse{}

	err = json.NewDecoder(resp.Body).Decode(&completionResp)

	if err != nil {
		return Message{}, err
	}

	return completionResp.Choices[0].Message, nil

}
