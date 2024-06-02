package main

import "fmt"

func init() {
	Cmd("clear", "Clear the terminal screen").WithHandler(
		func(c *command) {
			fmt.Print("\033[H\033[2J")
		},
	)
}
