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

func getBkt(b *bbolt.Bucket, key string) *bbolt.Bucket {
	s := b.Bucket([]byte(key))
	if s == nil {
		if hk, err := hex.DecodeString(key); err == nil {
			return b.Bucket(hk)
		}
	}
	return s
}

func getKey(b *bbolt.Bucket, key string) ([]byte, []byte, bool) {
	if getBkt(b, key) != nil {
		return nil, nil, false
	}
	hk := []byte(key)
	val := b.Get(hk)
	if val == nil {
		if k, err := hex.DecodeString(key); err == nil {
			val = b.Get(k)
			if val != nil {
				hk = k
			}
		}
	}
	return hk, val, true
}

func fuzzyMatch(pattern, subject string) bool {
	if pattern == "" {
		return true
	}
	idx := -1
	subject = strings.ToLower(subject)
	for _, p := range strings.ToLower(pattern) {
		x := strings.IndexRune(subject[idx+1:], p)
		if x == -1 {
			return false
		}
		idx += x + 1
	}
	return true
}

func pfxMatch(key []byte, pfx string) *prompt.Suggest {
	target := string(key)
	if !isPrintable(target) {
		target = hex.EncodeToString(key)
	}
	if fuzzyMatch(pfx, target) {
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

func hintKey(arg string) (ss []prompt.Suggest) {
	view(func(tx *bbolt.Tx) error {
		if b, _ := changeDir(tx); b != nil {
			b.ForEach(func(k, v []byte) error {
				if b.Bucket(k) != nil {
					return nil
				}
				if s := pfxMatch(k, arg); s != nil {
					ss = append(ss, *s)
				}
				return nil
			})
		}
		return nil
	})
	return
}

func hintBucket(arg string) (ss []prompt.Suggest) {
	view(func(tx *bbolt.Tx) error {
		b, err := changeDir(tx)
		if err != nil {
			return err
		}
		if b == nil {
			tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
				if s := pfxMatch(name, arg); s != nil {
					ss = append(ss, *s)
				}
				return nil
			})
			return nil
		}
		b.ForEach(func(k, v []byte) error {
			if b.Bucket(k) != nil {
				if s := pfxMatch(k, arg); s != nil {
					ss = append(ss, *s)
				}
			}
			return nil
		})
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
