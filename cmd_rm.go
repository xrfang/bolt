package main

import (
	"errors"
	"fmt"

	"go.etcd.io/bbolt"
)

func handleRm(c *command) {
	update(func(tx *bbolt.Tx) error {
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
		if key == "*" {
			cnt := countKeys(b)
			if cnt == 0 {
				fmt.Println("0 keys deleted.")
				return nil
			}
			hint := fmt.Sprintf("Are you sure to delete %d keys?", cnt)
			return confirmDo(hint, func() error {
				return b.ForEach(func(k, v []byte) error {
					return b.Delete(k)
				})
			})
		}
		return b.Delete([]byte(key))
	})
}

func init() {
	Cmd("rm", "Remove a key").WithParams("key").
		WithHandler(handleRm).WithCompleter(hintKey)
}
