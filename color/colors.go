package color

import "fmt"

var red string = "\033[31m"
var green string = "\033[32m"
var reset string = "\033[0m"
var yellow string = "\033[33m"

func Red(msg string) string {
	return fmt.Sprintf("%v%v%v", red, msg, reset)
}
func Green(msg string) string {
	return fmt.Sprintf("%v%v%v", green, msg, reset)
}
func Yellow(msg string) string {
	return fmt.Sprintf("%v%v%v", yellow, msg, reset)
}
