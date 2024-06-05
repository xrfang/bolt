package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/c-bata/go-prompt"
)

func init() {
	Cmd("help", "Show available commands").WithHandler(func(c *command) {
		var ss []prompt.Suggest
		cmds.Range(func(k, v any) bool {
			c := v.(*command)
			ss = append(ss, c.Suggest())
			return true
		})
		var maxlen int
		sort.Slice(ss, func(i, j int) bool {
			if l := len(ss[i].Text); l > maxlen {
				maxlen = l
			}
			return ss[i].Text < ss[j].Text
		})
		maxlen += 1
		for _, s := range ss {
			cmd := (s.Text + strings.Repeat(" ", maxlen))[:maxlen]
			fmt.Println(cmd, s.Description)
		}
	})
}
