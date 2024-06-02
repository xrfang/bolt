package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/fatih/color"
	"go.etcd.io/bbolt"
)

func ls(key, val []byte) {
	sfx := "/"
	if len(val) > 0 {
		sfx = ` (` + strconv.Itoa(len(val)) + `)`
	}
	if utf8.Valid(key) {
		fmt.Println(string(key) + sfx)
	} else {
		hl := color.New(color.FgHiRed)
		hl.Println(hex.EncodeToString(key) + sfx)
	}
}

func init() {
	Cmd("ls", "View keys under current bucket").WithHandler(
		func(c *command) {
			view(func(tx *bbolt.Tx) error {
				if len(bkt) < 2 {
					tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
						ls(name, nil)
						return nil
					})
					return nil
				}
				b, err := changeDir(tx)
				if err != nil {
					return err
				}
				b.ForEach(func(k, v []byte) error {
					ls(k, v)
					return nil
				})
				return nil
			})
		},
	)
}
