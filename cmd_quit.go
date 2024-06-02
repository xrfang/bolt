package main

import "os"

func init() {
	Cmd("exit", "Exit this tool").WithHandler(
		func(c *command) {
			resetTTY()
			os.Exit(0)
		},
	)
	Cmd("quit", "Exit this tool").WithHandler(
		func(c *command) {
			resetTTY()
			os.Exit(0)
		},
	)
}
