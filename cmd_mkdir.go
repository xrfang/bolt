package main

import (
	"errors"

	"go.etcd.io/bbolt"
)

func init() {
	Cmd("mkdir", "Create a new bucket").WithParams("dir").WithHandler(
		func(c *command) {
			update(func(tx *bbolt.Tx) error {
				dir := c.Arg("dir")
				if dir == "" {
					return errors.New("dir not specified")
				}
				b, err := changeDir(tx)
				if err != nil {
					return err
				}
				if b == nil {
					_, err = tx.CreateBucketIfNotExists([]byte(dir))
					return err
				}
				_, err = b.CreateBucketIfNotExists([]byte(dir))
				return err
			})
		},
	)
}
