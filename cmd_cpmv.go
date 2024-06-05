package main

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/c-bata/go-prompt"
	"go.etcd.io/bbolt"
)

func hintDest(arg string) (ss []prompt.Suggest) {
	view(func(tx *bbolt.Tx) error {
		mp, err := mergePath(arg)
		if err != nil {
			return err
		}
		b := tx.Bucket([]byte(mp[0]))
		if b == nil {
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

func getSrc(tx *bbolt.Tx, src string) ([]byte, []byte, error) {
	b, err := changeDir(tx)
	if err != nil {
		return nil, nil, err
	}
	notExist := fmt.Errorf("'%s' does not exist or is a bucket", src)
	if b == nil {
		return nil, nil, notExist
	}
	key, val, ok := getKey(b, src)
	if !ok {
		return nil, nil, notExist
	}
	return key, val, nil
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
	if len(base) == 0 {
		return nil, errors.New("invalid path: " + dst)
	}
	return base, nil
}

// TODO: multi是给src=*用的，移动所有key到某个目录
func getDst(tx *bbolt.Tx, dst string, multi bool) (*bbolt.Bucket, string, error) {
	mp, err := mergePath(dst)
	if err != nil {
		return nil, "", err
	}
	perr := errors.New("invalid path: " + dst)
	b := tx.Bucket([]byte(mp[0]))
	if b == nil {
		return nil, "", perr
	}
	for i, p := range mp[1:] {
		s := b.Bucket([]byte(p))
		if s == nil {
			if v := b.Get([]byte(p)); v != nil {
				return nil, "", fmt.Errorf("'%s' is an existing key", dst)
			}
			if i < len(mp)-2 {
				return nil, "", perr
			}
			return b, p, nil
		}
		b = s
	}
	return b, "", nil
}

func handleCp(c *command) {
	update(func(tx *bbolt.Tx) error {
		src := c.Arg("src")
		dst := c.Arg("dst")
		if dst == "" {
			return errors.New("missing destination")
		}
		b, bp, err := getDst(tx, dst, src == "*")
		if err != nil {
			return err
		}
		key, val, err := getSrc(tx, src)
		if err != nil {
			return err
		}
		if bp != "" {
			key = []byte(bp)
		}
		return b.Put(key, val)
	})
}

func handleMv(c *command) {
	update(func(tx *bbolt.Tx) error {
		src := c.Arg("src")
		dst := c.Arg("dst")
		if dst == "" {
			return errors.New("missing destination")
		}
		b, bp, err := getDst(tx, dst, src == "*")
		if err != nil {
			return err
		}
		key, val, err := getSrc(tx, src)
		if err != nil {
			return err
		}
		if bp != "" {
			key = []byte(bp)
		}
		if err := b.Put(key, val); err != nil {
			return err
		}
		if b, err = changeDir(tx); err == nil {
			err = b.Delete(key)
		}
		return err
	})
}

func init() {
	Cmd("cp", "Copy a key").WithParams("src", "dst").WithHandler(handleCp).
		WithCompleter(hintKey, hintDest)
	Cmd("mv", "Move (rename) a key").WithParams("src", "dst").WithHandler(handleMv).
		WithCompleter(hintKey, hintDest)
}
