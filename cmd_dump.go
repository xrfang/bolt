package main

import (
	"encoding/hex"
	"errors"
	"fmt"

	"go.etcd.io/bbolt"
)

func init() {
	Cmd("dump", "Hexdump the content of a key").WithParams("key").WithHandler(
		func(c *command) {
			view(func(tx *bbolt.Tx) error {
				key := c.Arg("key")
				if key == "" {
					return errors.New("key not specified")
				}
				b, err := changeDir(tx)
				if err != nil {
					return err
				}
				if b == nil {
					return fmt.Errorf("'%s' not exist or is a bucket", key)
				}
				val := b.Get([]byte(key))
				if len(val) == 0 {
					return fmt.Errorf("'%s' not exist or is a bucket", key)
				}
				fmt.Print(hex.Dump(val))
				return nil
			})
		},
	)
}
