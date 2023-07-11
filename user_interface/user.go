package userinterface

import (
	"errors"
	"fmt"

	"github.com/AntonyChR/terminalGPT/color"
)

func NewUserInterface() *UserInterfaceIO {
	return &UserInterfaceIO{
		PrintChannel: make(chan string),
	}
}

type UserInterfaceIO struct {
	PrintChannel chan string
}

func (u *UserInterfaceIO) GetInput(prompt string) (string, error) {
	var input string
	fmt.Print(color.Yellow(prompt))
	fmt.Scanln(&input)
	return input, validateInput(input)
}

func validateInput(input string) error {
	if input == "" || input == "\n" {
		return errors.New("invalid input")
	}
	return nil
}

func (u *UserInterfaceIO) PrintChannelMessages() {
	for {
		msg := <-u.PrintChannel
		fmt.Print(msg)
	}
}
