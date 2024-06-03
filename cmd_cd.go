package main

import (
	"github.com/c-bata/go-prompt"
	"go.etcd.io/bbolt"
)

func handleCd(c *command) {
	dir := c.Arg("dir")
	switch dir {
	case "":
		bkt = bkt[:1]
	case "..":
		if l := len(bkt); l > 1 {
			bkt = bkt[:l-1]
		}
	default:
		view(func(tx *bbolt.Tx) error {
			bkt = append(bkt, dir)
			_, err = changeDir(tx)
			if err != nil {
				bkt = bkt[:len(bkt)-1]
			}
			return err
		})
	}
}

func completeCd(c *command) (ss []prompt.Suggest) {
	view(func(tx *bbolt.Tx) error {
		pfx := c.Arg("dir")
		b, err := changeDir(tx)
		if err != nil {
			return err
		}
		if b == nil {
			tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
				if s := pfxMatch(name, pfx); s != nil {
					ss = append(ss, *s)
				}
				return nil
			})
			return nil
		}
		b.ForEach(func(k, v []byte) error {
			if len(v) == 0 {
				if s := pfxMatch(k, pfx); s != nil {
					ss = append(ss, *s)
				}
			}
			return nil
		})
		return nil
	})
	return
}

func init() {
	Cmd("cd", "Change current dir (open a bucket)").WithParams("dir").
		WithCompleter(completeCd).WithHandler(handleCd)
}
