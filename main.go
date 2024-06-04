package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/c-bata/go-prompt"
)

var readonly bool

func main() {
	defer resetTTY()
	flag.BoolVar(&readonly, "readonly", false, "read only mode")
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
	if flag.NArg() > 0 {
		fp, err := chkOpenArg(flag.Arg(0))
		if err == nil {
			err = open(fp, false)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	p := prompt.New(
		executor,
		completer,
		prompt.OptionLivePrefix(promptPrefix),
	)
	p.Run()
}
