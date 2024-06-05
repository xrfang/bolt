package main

import (
	"errors"

	"go.etcd.io/bbolt"
)

func handleAdd(c *command) {
	update(func(tx *bbolt.Tx) error {
		key, _ := decode(c.Arg("key"))
		if len(key) == 0 {
			return errors.New("missing or invalid key")
		}
		val, _ := decode(c.Arg("val"))
		if len(val) == 0 {
			return errors.New("missing or invalid value")
		}
		b, err := changeDir(tx)
		if err != nil {
			return err
		}
		if b == nil {
			return errors.New("cannot add key under root bucket")
		}
		if v := b.Get(key); v != nil {
			return errors.New("key already exists")
		}
		return b.Put(key, val)
	})
}

func init() {
	Cmd("add", "Add a new key").WithParams("key", "val").WithHandler(handleAdd)
}
