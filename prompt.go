package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"
	"go.xrfang.cn/shellwords"
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
	p := shellwords.NewParser()
	args, _ := p.Parse(text)
	if len(args) == 0 {
		return
	}
	cs, err := ParseCmd(args, false)
	if err != nil {
		fmt.Println("completer:", err.Error())
		return
	}
	switch len(cs) {
	case 0:
	case 1:
		if cs[0].name != args[0] {
			ss = append(ss, cs[0].Suggest())
		} else if len(args) > 1 || strings.HasSuffix(text, " ") {
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
	p := shellwords.NewParser()
	args, err := p.Parse(cmdline)
	if len(args) == 0 {
		if err != nil {
			fmt.Println("ERROR:", err)
		}
		return
	}
	cs, err := ParseCmd(args, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	cs[0].Run()
}
