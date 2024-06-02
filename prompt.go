package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"
)

func resetTTY() {
	cmd := exec.Command("/bin/stty", "-raw", "echo")
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func promptPrefix() (string, bool) {
	if len(bkt) == 0 {
		return "(none)> ", true
	}
	pfx := filepath.Base(bkt[0])
	if len(bkt) > 1 {
		pfx += ":" + strings.Join(bkt[1:], "/")
	}
	return fmt.Sprintf("%s> ", pfx), true
}

func completer(d prompt.Document) (ss []prompt.Suggest) {
	// text := d.TextBeforeCursor()
	// p := shellwords.NewParser()
	// args, err := p.Parse(text)
	// if err != nil {
	// fmt.Println(err.Error())
	// return
	// }
	// fmt.Println("completer:", len(args))
	return
}

func executor(cmdline string) {
	cs, err := ParseCmd(cmdline)
	if err != nil {
		fmt.Println(err)
		return
	}
	cs[0].Run()
}
