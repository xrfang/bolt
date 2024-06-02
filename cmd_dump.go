package main

import (
	"encoding/hex"
	"fmt"

	"go.etcd.io/bbolt"
)

func init() {
	Cmd("dump", "Hexdump the content of a key").WithParams("key").WithHandler(
		func(c *command) {
			key := c.Arg("key")
			view(func(tx *bbolt.Tx) error {
				if len(bkt) < 2 {
					fmt.Printf("'%s' not exist or is a bucket\n", key)
					return nil
				}
				b, err := changeDir(tx)
				if err != nil {
					return err
				}
				val := b.Get([]byte(key))
				if len(val) == 0 {
					fmt.Printf("'%s' not exist or is a bucket\n", key)
				} else {
					fmt.Print(hex.Dump(val))
				}
				return nil
			})
		},
	)
}
