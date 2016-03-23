package ui

import (
	"fmt"
	"strings"
)

type Runtime string

func (r Runtime) String() string {
	switch strings.ToLower(string(r)) {
	case strings.ToLower(DEA):
		return DEA
	case strings.ToLower(Diego):
		return Diego
	default:
		return string(r)
	}
}

const (
	DEA   = "DEA"
	Diego = "Diego"
)

func sayOk() {
	fmt.Println(say("Ok\n", 32, 1))
}

func say(message string, color uint, bold int) string {
	return fmt.Sprintf("\033[%d;%dm%s\033[0m", bold, color, message)
}

type ApplicationPrinter interface {
	Name() string
	Organization() string
	Space() string
}
