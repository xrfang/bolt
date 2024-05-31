package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"
)

func main() {
	defer func() {
		cmd := exec.Command("/bin/stty", "-raw", "echo")
		cmd.Stdin = os.Stdin
		cmd.Run()
	}()
	ver := flag.Bool("version", false, "show version info")
	flag.Usage = func() {
		fmt.Println("BoltDB Editor", verinfo())
		fmt.Printf("\nUSAGE: %s [OPTIONS] [db file]\n", filepath.Base(os.Args[0]))
		fmt.Println("\nOPTIONS:")
		flag.PrintDefaults()
	}
	flag.Parse()
	if *ver {
		fmt.Println(verinfo())
		return
	}
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefixTextColor(prompt.Cyan),
		prompt.OptionLivePrefix(promptPrefix),
		prompt.OptionSetExitCheckerOnInput(func(in string, breakline bool) bool {
			if breakline {
				switch strings.ToLower(in) {
				case "exit", "quit":
					return true
				}
			}
			return false
		}),
	)
	p.Run()
}
