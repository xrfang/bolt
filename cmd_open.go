package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"go.etcd.io/bbolt"
)

func connect(fp string) error {
	mode := fs.FileMode(0666)
	if readonly {
		mode = fs.FileMode(0444)
	}
	db, err = bbolt.Open(fp, mode, nil)
	if err == nil {
		bkt = []string{fp}
	}
	return err
}

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
	return connect(fp)
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

func init() {
	Cmd("create", "Create a new BoltDB database").WithParams("filename").WithHandler(
		func(c *command) {
			fp, err := chkOpenArg(c.Arg("filename"))
			if err != nil {
				fmt.Println(err)
				return
			}
			if err := open(fp, true); err != nil {
				fmt.Println(err)
			}
		},
	)
	Cmd("open", "Open a BoltDB database").WithParams("filename").WithHandler(
		func(c *command) {
			fp, err := chkOpenArg(c.Arg("filename"))
			if err != nil {
				fmt.Println(err)
				return
			}
			if err := open(fp, true); err != nil {
				fmt.Println(err)
			}
		},
	)
}
