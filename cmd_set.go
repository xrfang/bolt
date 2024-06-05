package main

import (
	"errors"
	"fmt"

	"go.etcd.io/bbolt"
)

func handleSet(c *command) {
	update(func(tx *bbolt.Tx) error {
		key := c.Arg("key")
		val, _ := decode(c.Arg("val"))
		if len(val) == 0 {
			return errors.New("missing or invalid value")
		}
		b, err := changeDir(tx)
		if err != nil {
			return err
		}
		if b == nil {
			return fmt.Errorf("'%s' does not exist or is a bucket", key)
		}
		k, v, ok := getKey(b, key)
		if !ok || v == nil {
			return fmt.Errorf("'%s' does not exist or is a bucket", key)
		}
		return b.Put(k, val)
	})
}

func init() {
	Cmd("set", "Set the value of a key").WithParams("key", "val").
		WithHandler(handleSet).WithCompleter(hintKey)
}
