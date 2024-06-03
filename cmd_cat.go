package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
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
			return fmt.Errorf("'%s' not exist or is a bucket", key)
		}
		val := b.Get([]byte(key))
		if len(val) == 0 {
			return fmt.Errorf("'%s' not exist or is a bucket", key)
		}
		if isPrintable(string(val)) {
			fmt.Println(string(val))
		} else {
			hl := color.New(color.FgHiRed)
			hl.Println(hex.EncodeToString(val))
		}
		return nil
	})
}

func completeCat(c *command) (ss []prompt.Suggest) {
	view(func(tx *bbolt.Tx) error {
		if b, _ := changeDir(tx); b != nil {
			pfx := c.Arg("key")
			b.ForEach(func(k, v []byte) error {
				if len(v) == 0 {
					return nil
				}
				key := string(k)
				if !isPrintable(key) {
					key = hex.EncodeToString(k)
				}
				if pfx != "" {
					fmt.Printf("key=%s; pfx=%s\n", key, pfx)
				}
				if pfx == "" || strings.HasPrefix(key, pfx) {
					ss = append(ss, prompt.Suggest{Text: key})
				}
				return nil
			})
		}
		return nil
	})
	return
}

func init() {
	Cmd("cat", "Show content of a key").WithParams("key").
		WithHandler(handleCat).WithCompleter(completeCat)
}
