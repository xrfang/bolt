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
	text := d.CurrentLineBeforeCursor() // CurrentLine() // d.TextBeforeCursor()
	if text == "" {
		return
	}
	cs, err := ParseCmd(text, false)
	if err != nil {
		fmt.Println("completer:", err.Error())
		return
	}
	switch len(cs) {
	case 0:
	case 1:
		//TODO: 这里有问题，需要检查是否完整匹配！！！
		if strings.HasSuffix(text, " ") {
			ss = cs[0].SuggestNextArg()
		}
	default:
		for _, c := range cs {
			ss = append(ss, c.Suggest())
		}
	}
	return
}

func executor(cmdline string) {
	cs, err := ParseCmd(cmdline, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	cs[0].Run()
}
