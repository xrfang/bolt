package main

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/fatih/color"
	"go.etcd.io/bbolt"
)

func handleCat(c *command) {
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
			return fmt.Errorf("'%s' does not exist or is a bucket", key)
		}
		_, val, ok := getKey(b, key)
		if !ok {
			return fmt.Errorf("'%s' does not exist or is a bucket", key)
		}
		if len(val) > 0 {
			if isPrintable(string(val)) {
				fmt.Println(string(val))
			} else {
				hl := color.New(color.FgHiRed)
				hl.Println(hex.EncodeToString(val))
			}
		}
		return nil
	})
}

func init() {
	Cmd("cat", "Show content of a key").WithParams("key").
		WithHandler(handleCat).WithCompleter(hintKey)
}
