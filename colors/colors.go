package color

import "fmt"

var COLORS = map[string]string{
	"red":    "\033[31m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
	"reset":  "\033[0m",
}

func Colorize(color string, msg string) string {
	return fmt.Sprintf("%v%v%v", COLORS[color], msg, COLORS["reset"])
}

func Red(msg string) string {
	return Colorize("red", msg)
}

func Green(msg string) string {
	return Colorize("green", msg)
}

func Yellow(msg string) string {
	return Colorize("yellow", msg)
}
