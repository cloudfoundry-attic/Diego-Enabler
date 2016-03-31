package ui

import "github.com/fatih/color"

func SayOK() {
	c := color.New(color.FgGreen).Add(color.Bold)
	c.Println("OK\n")
}

func SayFailed() {
	c := color.New(color.FgRed).Add(color.Bold)
	c.Println("FAILED")
}
