package main

import (
	"path"
	"strings"

	"go.etcd.io/bbolt"
)

func handleCd(c *command) {
	dir := c.Arg("dir")
	if dir == "" {
		bkt = bkt[:1]
		return
	}
	dir = path.Clean(dir)
	if dir == ".." {
		if l := len(bkt); l > 1 {
			bkt = bkt[:l-1]
		}
		return
	}
	view(func(tx *bbolt.Tx) error {
		old := append([]string{}, bkt...)
		if strings.HasPrefix(dir, "/") {
			dir = dir[1:]
		} else {
			dir = path.Join(strings.Join(bkt[1:], "/"), dir)
		}
		bkt = append([]string{bkt[0]}, strings.Split(path.Clean(dir), "/")...)
		_, err = changeDir(tx)
		if err != nil {
			bkt = old
		}
		return err
	})
}

func init() {
	Cmd("cd", "Change current dir (open a bucket)").WithParams("dir").
		WithCompleter(hintPath).WithHandler(handleCd)
}
