package displayhelpers

import "fmt"

func Say(message string, color uint, bold int) string {
	return fmt.Sprintf("\033[%d;%dm%s\033[0m", bold, color, message)
}

func SayOK() {
	fmt.Println(Say("Ok\n", 32, 1))
}

func SayFailed() {
	fmt.Println(Say("FAILED", 31, 1))
}
