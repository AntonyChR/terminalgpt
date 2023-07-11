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

	"github.com/AntonyChR/terminalGPT/color"
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
	StreamData  bool    //defaul false
}

func NewChat(c Configuration) *Openai {
	initialMessage := Message{
		Role:    ChatRoles.System,
		Content: "Eres un asistente, si tu respuesta contiene código de programación los nombres de variable, funciones, etc... deben ser en ingles. El resto de la respuesta debe estar en el mismo idioma en el que se hizo la pregunta",
	}
	return &Openai{
		apikey:           c.Apikey,
		url:              c.ApiUrl,
		IncommingMessage: make(chan string),
		chat: Chat{
			Model: c.Model,
			Messages: []Message{
				initialMessage,
			},
			Temperature: c.Temperature,
			Stream:      c.StreamData,
		},
	}
}

type Openai struct {
	apikey           string
	url              string
	chat             Chat
	IncommingMessage chan string
}

type Chat struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Temperature float32   `json:"temperature,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
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
	req.Header.Set("Authorization", "Bearer "+o.apikey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Message{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Message{}, &customErrors.RequestError{StatusCode: resp.StatusCode, Err: errors.New("unespected request errror")}
	}

	if !o.chat.Stream {

		completionResp := CompletionResponse{}

		err = json.NewDecoder(resp.Body).Decode(&completionResp)

		if err != nil {
			return Message{}, err
		}

		o.IncommingMessage <- completionResp.Choices[0].Message.Content
		o.IncommingMessage <- "\n"

		return completionResp.Choices[0].Message, nil
	}

	chunkBuffer := make([]byte, 4096)
	re := regexp.MustCompile(`(?s)"finish_reason":(null|"stop").*`)
	msg := Message{Role: ChatRoles.Assistant}
	content := ""
	for {
		chunkSize, err := resp.Body.Read(chunkBuffer)
		fmt.Println(resp.Header.Get("Content-Type"))		
		if err != nil {
			if errors.Is(err, io.EOF) {
				o.IncommingMessage <- "\n"
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
				fmt.Println(color.Red(err.Error()))
				continue
			}
			deltaString := chunkObject.Choices[0].Delta.Content
			content += deltaString
			o.IncommingMessage <- deltaString
		}

	}

	msg.Content = content
	o.IncommingMessage <- "\n[A]: "

	return msg, nil

}

func (o *Openai) Reset() {
	fmt.Println(color.Red("Reset context"))
	o.chat.Messages = []Message{o.chat.Messages[0]}
}

func (o *Openai) ListenAndPrintIncommingMsg() {
	for {
		msg := <-o.IncommingMessage
		fmt.Print(color.Green(msg))
	}
}
