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
	text := d.CurrentLineBeforeCursor()
	argv, cs := ParseCmd(text)
	switch len(cs) {
	case 0:
	case 1:
		if cs[0].name != argv[0] {
			ss = append(ss, cs[0].Suggest())
		} else if len(argv) > 1 || strings.HasSuffix(text, " ") {
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
	argv, cs := ParseCmd(cmdline)
	if len(cs) != 1 || cs[0].name != argv[0] {
		fmt.Printf("unknown command '%s' (try 'help')\n", argv[0])
		return
	}
	cs[0].Run()
}
