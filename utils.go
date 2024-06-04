package main

import (
	"encoding/hex"
	"strings"
	"unicode"

	"github.com/c-bata/go-prompt"
	"go.etcd.io/bbolt"
)

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

func getKey(b *bbolt.Bucket, key string) []byte {
	val := b.Get([]byte(key))
	if len(val) == 0 {
		if hk, err := hex.DecodeString(key); err == nil {
			return b.Get(hk)
		}
	}
	return val
}

func pfxMatch(key []byte, pfx string) *prompt.Suggest {
	target := string(key)
	if !isPrintable(target) {
		target = hex.EncodeToString(key)
	}
	if pfx == "" || strings.HasPrefix(strings.ToUpper(target),
		strings.ToUpper(pfx)) {
		return &prompt.Suggest{Text: target}
	}
	return nil
}

func countKeys(b *bbolt.Bucket) (cnt int) {
	b.ForEach(func(k, v []byte) error {
		if len(v) > 0 {
			cnt++
		}
		return nil
	})
	return
}

func confirmDo(hint string, f func() error) error {
	var err error
	var done bool
	p := prompt.New(
		func(cmd string) {
			switch cmd {
			case "yes":
				err = f()
				fallthrough
			case "no":
				done = true
			}
		},
		func(d prompt.Document) []prompt.Suggest {
			return []prompt.Suggest{{Text: "yes"}, {Text: "no"}}
		},
		prompt.OptionPrefix(hint+" "),
		prompt.OptionSetExitCheckerOnInput(func(in string, breakline bool) bool {
			return breakline && (in == "yes" || in == "no")
		}),
		prompt.OptionLivePrefix(func() (prefix string, useLivePrefix bool) {
			if done {
				return "", true
			}
			return hint + " ", true
		}),
	)
	p.Run()
	return err
}
