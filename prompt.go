package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"go.xrfang.cn/shellwords"
)

// 自定义提示符函数
func promptPrefix() (string, bool) {
	dir, err := os.Getwd()
	if err != nil {
		return "> ", true
	}
	return fmt.Sprintf("%s> ", dir), true
}

func completer(d prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "exit", Description: "Exit this tool"},
		{Text: "ls", Description: "View keys under current bucket"},
		{Text: "open", Description: "Open a BoltDB database"},
		{Text: "quit", Description: "Exit this tool"},
	}
}

// 交互模式命令执行函数
func executor(cmd string) {
	p := shellwords.NewParser()
	args, err := p.Parse(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if len(args) == 0 {
		return
	}
	cmd = strings.ToLower(args[0])
	switch cmd {
	case "quit", "exit": //直接退出
	case "open":
		fmt.Printf("TODO: %+v\n", args)
	case "ls":
		fmt.Printf("TODO: %+v\n", args)
	default:
		fmt.Printf("unknown command '%s'\n", cmd)
	}
}
