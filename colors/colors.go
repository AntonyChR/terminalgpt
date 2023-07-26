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
