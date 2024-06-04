package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/c-bata/go-prompt"
)

func open(fp string, create bool) (err error) {
	st, err := os.Stat(fp)
	if err == nil && st.IsDir() {
		return fmt.Errorf("'%s' is a directory (file expected)", fp)
	}
	if create {
		if err == nil {
			return fmt.Errorf("'%s' already exists", fp)
		}
	} else if err != nil {
		return fmt.Errorf("'%s' does not exist", fp)
	}
	return openBoltDB(fp)
}

func chkOpenArg(fn string) (fp string, err error) {
	if fn == "" {
		return "", errors.New("missing file name")
	}
	bkt = []string{}
	if db != nil {
		db.Close()
	}
	return filepath.Abs(fn)
}

func handleCreate(c *command) {
	fp, err := chkOpenArg(c.Arg("filename"))
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := open(fp, true); err != nil {
		fmt.Println(err)
	}
}

func handleOpen(c *command) {
	fp, err := chkOpenArg(c.Arg("filename"))
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := open(fp, false); err != nil {
		fmt.Println(err)
	}
}

func hintFile(arg string) (ss []prompt.Suggest) {
	fns, _ := filepath.Glob("*")
	for _, fn := range fns {
		if fuzzyMatch(arg, fn) {
			ss = append(ss, prompt.Suggest{Text: fn})
		}
	}
	return
}

func init() {
	Cmd("create", "Create a new BoltDB database").WithParams("filename").
		WithHandler(handleCreate)
	Cmd("open", "Open a BoltDB database").WithParams("filename").
		WithHandler(handleOpen).WithCompleter(hintFile)
}
