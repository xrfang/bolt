package main

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"path"
	"strings"
	"unicode"

	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"
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

func hintKey(arg string) (ss []prompt.Suggest) {
	hint(func(tx *bbolt.Tx) error {
		if b, _ := changeDir(tx); b != nil {
			b.ForEach(func(k, v []byte) error {
				if b.Bucket(k) != nil {
					return nil
				}
				if s := hintMatch(k, arg); s != nil {
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
	hint(func(tx *bbolt.Tx) error {
		b, err := changeDir(tx)
		if err != nil {
			return err
		}
		if b == nil {
			tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
				if s := hintMatch(name, arg); s != nil {
					ss = append(ss, *s)
				}
				return nil
			})
			return nil
		}
		b.ForEach(func(k, v []byte) error {
			if b.Bucket(k) != nil {
				if s := hintMatch(k, arg); s != nil {
					ss = append(ss, *s)
				}
			}
			return nil
		})
		return nil
	})
	return
}

func mergePath(dst string) ([]string, error) {
	dst = path.Clean(dst)
	if strings.HasPrefix(dst, "/") {
		return strings.Split(dst[1:], "/"), nil
	}
	var base []string
	if len(bkt) > 1 {
		base = append(base, bkt[1:]...)
	}
	for _, d := range strings.Split(dst, "/") {
		switch d {
		case ".":
		case "..":
			if len(base) == 0 {
				return nil, errors.New("invalid path: " + dst)
			}
			base = base[:len(base)-1]
		default:
			base = append(base, d)
		}
	}
	return base, nil
}

func hintPath(arg string) (ss []prompt.Suggest) {
	hint(func(tx *bbolt.Tx) error {
		mp, err := mergePath(arg)
		if err != nil {
			return err
		}
		if len(mp) == 0 {
			tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
				ss = append(ss, prompt.Suggest{Text: string(name)})
				return nil
			})
			return nil
		}
		b := tx.Bucket([]byte(mp[0]))
		if b == nil {
			tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
				if fuzzyMatch(mp[0], string(name)) {
					ss = append(ss, prompt.Suggest{Text: string(name)})
				}
				return nil
			})
			return nil
		}
		var pfx string
		for _, p := range mp[1:] {
			s := b.Bucket([]byte(p))
			if s == nil {
				pfx = p
				break
			}
			b = s
		}
		if pfx != "" {
			mp = mp[:len(mp)-1]
		} else {
			ss = []prompt.Suggest{{Text: "/" + strings.Join(mp, "/")}}
		}
		b.ForEachBucket(func(k []byte) error {
			var hp string
			if pfx == "" || strings.HasPrefix(string(k), pfx) {
				hp = strings.Join(append(mp, string(k)), "/")
				ss = append(ss, prompt.Suggest{Text: "/" + hp})
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

// 该函数处理set和add的value参数，“\x”开头表示16进制，“\b”开头表示base64编码，
// 除以上两种情形外，如果value本身以“\\”开头，将其转换为“\”。
func decode(val string) ([]byte, error) {
	switch {
	case strings.HasPrefix(val, `\x`):
		return hex.DecodeString(val[2:])
	case strings.HasPrefix(val, `\b`):
		val = strings.ReplaceAll(strings.ReplaceAll(val[2:],
			"-", "+"), "_", "/")
		switch len(val) % 4 {
		case 2:
			val += "=="
		case 3:
			val += "="
		}
		return base64.StdEncoding.DecodeString(val)
	case strings.HasPrefix(val, `\\`):
		return []byte(val[1:]), nil
	default:
		return []byte(val), nil
	}
}

// 该函数用连续的空格切分arg为最多n份
func pargs(arg string, n int) (args []string) {
	const (
		backSlash = "\xFE\x01\xFF"
		space     = "\xFE\x02\xFF"
		newLine   = "\xFE\x03\xFF"
	)
	if arg = strings.TrimSpace(arg); arg == "" || n < 1 {
		return nil
	}
	arg = strings.ReplaceAll(arg, `\\`, backSlash)
	arg = strings.ReplaceAll(arg, `\ `, space)
	arg = strings.ReplaceAll(arg, `\n`, newLine)
	defer func() {
		for i := range args {
			args[i] = strings.ReplaceAll(args[i], backSlash, `\`)
			args[i] = strings.ReplaceAll(args[i], space, ` `)
			args[i] = strings.ReplaceAll(args[i], newLine, "\n")
		}
	}()
	var word []byte
	for i, b := range []byte(arg) {
		if b != ' ' {
			if len(args) == n-1 {
				word = append(word, []byte(arg[i:])...)
				args = append(args, string(word))
				return
			}
			word = append(word, b)
		} else if len(word) > 0 {
			args = append(args, string(word))
			word = nil
		}
	}
	if len(word) > 0 {
		args = append(args, string(word))
	}
	return
}

var summary func(string, ...any)

func init() {
	summary = color.New(color.FgHiBlack).PrintfFunc()
}
