package main

import (
	"strings"
	"sync"

	"github.com/c-bata/go-prompt"
)

type (
	hintFunc func(string) []prompt.Suggest
	command  struct {
		name string
		desc string
		para []string
		exec func(*command)
		cmpl []hintFunc //func(*command) []prompt.Suggest
		_arg map[string]string
	}
)

func (c *command) WithParams(para ...string) *command {
	c.para = para
	return c
}

func (c *command) WithHandler(f func(*command)) *command {
	c.exec = f
	return c
}

func (c *command) WithCompleter(hfs ...hintFunc) *command {
	c.cmpl = hfs
	return c
}

func (c *command) withArgs(args []string) *command {
	c._arg = make(map[string]string)
	for i, p := range c.para {
		if i < len(args) {
			c._arg[p] = args[i]
		}
	}
	return c
}

func (c *command) Arg(name string) string {
	return c._arg[name]
}

func (c *command) Suggest() prompt.Suggest {
	return prompt.Suggest{Text: c.name, Description: c.desc}
}

func (c *command) SuggestNextArg() []prompt.Suggest {
	la := len(c._arg) - 1
	var arg string
	if la >= 0 {
		arg = c._arg[c.para[la]]
	} else {
		la = 0
	}
	if len(c.cmpl) <= la {
		return nil
	}
	return c.cmpl[la](arg)
}

func (c *command) Run() {
	c.exec(c)
}

func Cmd(name, desc string) *command {
	c := &command{name: name, desc: desc, _arg: make(map[string]string)}
	cmds.Store(c.name, c)
	return c
}

func ParseCmd(cmdline string) (argv []string, cs []*command) {
	ca := strings.SplitN(strings.TrimSpace(cmdline), " ", 2)
	cmd := ca[0]
	if cmd == "" {
		return
	}
	v, _ := cmds.Load(cmd)
	if v != nil {
		argv = append(argv, cmd)
		c := v.(*command)
		if len(ca) > 1 {
			args := pargs(ca[1], len(c.para))
			c.withArgs(args)
			argv = append(argv, args...)
		}
		cs = append(cs, c)
		return
	}
	cmds.Range(func(k, v any) bool {
		if strings.HasPrefix(k.(string), cmd) {
			argv = append(argv, cmd)
			c := v.(*command)
			if len(ca) > 1 {
				args := pargs(ca[1], len(c.para))
				c.withArgs(args)
				argv = append(argv, args...)
			}
			cs = append(cs, c)
		}
		return true
	})
	return
}

var cmds sync.Map //map[string]*command
