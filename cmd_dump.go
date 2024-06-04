package main

import (
	"encoding/hex"
	"errors"
	"fmt"

	"go.etcd.io/bbolt"
)

func handleDump(c *command) {
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
		fmt.Print(hex.Dump(val))
		return nil
	})
}

func init() {
	Cmd("dump", "Hexdump the content of a key").WithParams("key").
		WithCompleter(hintKey).WithHandler(handleDump)
}
