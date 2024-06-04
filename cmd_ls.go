package main

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"go.etcd.io/bbolt"
)

func ls(key, val []byte) {
	sfx := "/"
	if val != nil {
		sfx = ` (` + strconv.Itoa(len(val)) + ` bytes)`
	}
	if isPrintable(string(key)) {
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
					if b.Bucket(k) != nil {
						v = nil
					} else if v == nil {
						v = []byte{}
					}
					ls(k, v)
					return nil
				})
				return nil
			})
		},
	)
}
