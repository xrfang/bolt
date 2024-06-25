package main

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"go.etcd.io/bbolt"
)

func ls(key, val []byte) {
	sfx := "/"
	if val != nil {
		sfx = color.HiBlackString(" (%d bytes)", len(val))
	}
	if isPrintable(string(key)) {
		fmt.Println(string(key) + sfx)
	} else {
		fmt.Println(color.HiRedString(hex.EncodeToString(key)) + sfx)
	}
}

func handlLs(c *command) {
	var dir func(b *bbolt.Bucket, key string) int
	dir = func(b *bbolt.Bucket, key string) int {
		cnt := 0
		b.ForEach(func(k, v []byte) error {
			if strings.HasSuffix(key, "/") {
				if string(k) == key[:len(key)-1] {
					if s := b.Bucket(k); s != nil {
						cnt = dir(s, "")
					}
				}
				return nil
			}
			if key != "" && !wildMatch(k, key, true) {
				return nil
			}
			if b.Bucket(k) != nil {
				v = nil
			} else if v == nil {
				v = []byte{}
			}
			cnt++
			ls(k, v)
			return nil
		})
		return cnt
	}
	view(func(tx *bbolt.Tx) (err error) {
		cnt := 0
		key := c.Arg("key")
		defer func() {
			if err == nil && cnt > 0 {
				summary("%d items\n", cnt)
			}
		}()
		if len(bkt) < 2 {
			tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
				if strings.HasSuffix(key, "/") {
					if string(name) == key[:len(key)-1] {
						if b := tx.Bucket(name); b != nil {
							cnt = dir(b, "")
						}
					}
					return nil
				}
				if key == "" || wildMatch(name, key, true) {
					cnt++
					ls(name, nil)
				}
				return nil
			})
			return nil
		}
		b, err := changeDir(tx)
		if err != nil {
			return err
		}
		cnt = dir(b, key)
		return nil
	})
}

func init() {
	Cmd("ls", "View keys under current bucket").WithParams("key").
		WithHandler(handlLs).WithCompleter(hintAny)
}
