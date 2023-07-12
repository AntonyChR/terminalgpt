package userinterface

import (
	"bufio"
	"errors"
	"fmt"
	"os"

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
	fmt.Print(color.Yellow(prompt))
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

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
