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
		hk, val, ok := getKey(b, key)
		if !ok {
			return fmt.Errorf("'%s' does not exist or is a bucket", key)
		}
		if val != nil {
			return b.Delete(hk)
		}
		var todo [][]byte
		b.ForEach(func(k, v []byte) error {
			if b.Bucket(k) == nil {
				if wildMatch(k, key) {
					todo = append(todo, k)
				}
			}
			return nil
		})
		if len(todo) == 0 {
			fmt.Println("0 keys deleted")
			return nil
		}
		var hint string
		if len(todo) == 1 {
			hint = fmt.Sprintf("Delete '%s'?", todo[0])
		} else {
			hint = fmt.Sprintf("Are you sure to delete %d keys?", len(todo))
		}
		return confirmDo(hint, func() error {
			for _, k := range todo {
				if err := b.Delete(k); err != nil {
					return err
				}
			}
			fmt.Printf("deleted %d keys\n", len(todo))
			return nil
		})
	})
}

func init() {
	Cmd("rm", "Remove a key").WithParams("key").
		WithHandler(handleRm).WithCompleter(hintKey)
}
