package main

import (
	"encoding/hex"
	"fmt"

	"github.com/fatih/color"
	"go.etcd.io/bbolt"
)

func ls(key, val []byte) {
	sfx := "/"
	if val != nil {
		sfx = color.HiBlackString(" (%d bytes)", len(val))
	}
	if isPrintable(string(key)) {
		fmt.Println(string(key) + sfx)
	} else {
		fmt.Println(color.HiRedString(hex.EncodeToString(key)) + sfx)
	}
}

func init() {
	Cmd("ls", "View keys under current bucket").WithParams("key").WithHandler(
		func(c *command) {
			view(func(tx *bbolt.Tx) (err error) {
				cnt := 0
				key := c.Arg("key")
				defer func() {
					if err == nil && cnt > 0 {
						summary("%d items\n", cnt)
					}
				}()
				if len(bkt) < 2 {
					tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
						if key == "" || wildMatch(name, key, true) {
							cnt++
							ls(name, nil)
						}
						return nil
					})
					return nil
				}
				b, err := changeDir(tx)
				if err != nil {
					return err
				}
				b.ForEach(func(k, v []byte) error {
					if key != "" && !wildMatch(k, key, true) {
						return nil
					}
					if b.Bucket(k) != nil {
						v = nil
					} else if v == nil {
						v = []byte{}
					}
					cnt++
					ls(k, v)
					return nil
				})
				return nil
			})
		},
	)
}
