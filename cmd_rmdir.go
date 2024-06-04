package main

import (
	"errors"
	"fmt"

	"go.etcd.io/bbolt"
)

func handleRmdir(c *command) {
	update(func(tx *bbolt.Tx) error {
		dir := c.Arg("dir")
		if dir == "" {
			return errors.New("dir not specified")
		}
		b, err := changeDir(tx)
		if err != nil {
			return err
		}
		var sb *bbolt.Bucket
		if b == nil {
			sb = tx.Bucket([]byte(dir))
		} else {
			sb = b.Bucket([]byte(dir))
		}
		if sb == nil {
			return fmt.Errorf("'%s' not exist", dir)
		}
		if sb.ForEach(func(k, v []byte) error { return errors.New("has key") }) != nil {
			return fmt.Errorf("'%s' is not empty", dir)
		}
		if b == nil {
			return tx.DeleteBucket([]byte(dir))
		}
		return b.DeleteBucket([]byte(dir))
	})
}

func init() {
	Cmd("rmdir", "Remove a bucket").WithParams("dir").WithHandler(handleRmdir).
		WithCompleter(completeBucket)
}
