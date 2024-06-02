package main

import (
	"encoding/hex"
	"fmt"
	"unicode/utf8"

	"github.com/fatih/color"
	"go.etcd.io/bbolt"
)

func cat(val []byte) {
	if utf8.Valid(val) {
		fmt.Println(string(val))
	} else {
		hl := color.New(color.FgHiRed)
		hl.Println(hex.EncodeToString(val))
	}
}

func init() {
	Cmd("cat", "Show content of a key").WithParams("key").WithHandler(
		func(c *command) {
			view(func(tx *bbolt.Tx) error {
				key := c.Arg("key")
				if len(bkt) < 2 {
					fmt.Printf("'%s' not exist or is a bucket\n", key)
					return nil
				}
				b, err := changeDir(tx)
				if err != nil {
					return err
				}
				val := b.Get([]byte(key))
				cat(val)
				return nil
			})
		},
	)
}
