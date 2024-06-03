package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/fatih/color"
	"go.etcd.io/bbolt"
)

func init() {
	Cmd("cat", "Show content of a key").WithParams("key").WithHandler(
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
				if utf8.Valid(val) {
					fmt.Println(string(val))
				} else {
					hl := color.New(color.FgHiRed)
					hl.Println(hex.EncodeToString(val))
				}
				return nil
			})
		},
	)
}
