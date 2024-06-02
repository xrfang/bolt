package main

import (
	"errors"
	"strings"
	"sync"

	"go.xrfang.cn/shellwords"
)

type (
	command struct {
		name string
		desc string
		para []string
		exec func(*command)
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

func (c *command) WithArgs(args []string) *command {
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

func (c *command) Run() {
	c.exec(c)
}

func Cmd(name, desc string) *command {
	c := &command{name: name, desc: desc, _arg: make(map[string]string)}
	cmds.Store(c.name, c)
	return c
}

func ParseCmd(cmdline string) (cs []*command, err error) {
	p := shellwords.NewParser()
	args, err := p.Parse(cmdline)
	if err != nil {
		return nil, err
	}
	if len(args) == 0 {
		return nil, nil
	}
	v, _ := cmds.Load(args[0])
	if v != nil {
		c := v.(*command)
		cs = append(cs, c.WithArgs(args[1:]))
		return
	}
	cmds.Range(func(k, v any) bool {
		if strings.HasPrefix(k.(string), args[0]) {
			c := v.(*command)
			cs = append(cs, c.WithArgs(args[1:]))
		}
		return true
	})
	if len(cs) != 1 {
		err = errors.New("unknown command: " + args[0])
	}
	return
}

var cmds sync.Map //map[string]*command
