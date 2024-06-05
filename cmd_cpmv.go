package main

import (
	"errors"
	"fmt"

	"go.etcd.io/bbolt"
)

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
	if !ok || val == nil {
		return nil, nil, notExist
	}
	return key, val, nil
}

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
			if multi { //source is multiple keys, destination must be a bucket
				return nil, "", perr
			}
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

func batchOp(tx *bbolt.Tx, dst *bbolt.Bucket, del bool) (err error) {
	src, err := changeDir(tx)
	if err != nil {
		return err
	}
	return src.ForEach(func(k, v []byte) error {
		if src.Bucket(k) != nil {
			return nil
		}
		if err := dst.Put(k, v); err != nil {
			return err
		}
		if del {
			return src.Delete(k)
		}
		return nil
	})
}

func handleCp(c *command) {
	update(func(tx *bbolt.Tx) error {
		src := c.Arg("src")
		dst := c.Arg("dst")
		if dst == "" {
			return errors.New("missing destination")
		}
		var multi bool
		key, val, err := getSrc(tx, src)
		if err != nil {
			if src != "*" {
				return err
			}
			multi = true
		}
		b, bp, err := getDst(tx, dst, multi)
		if err != nil {
			return err
		}
		if bp != "" {
			key = []byte(bp)
		}
		if multi {
			return batchOp(tx, b, false)
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
		var multi bool
		key, val, err := getSrc(tx, src)
		if err != nil {
			if src != "*" {
				return err
			}
			multi = true
		}
		b, bp, err := getDst(tx, dst, multi)
		if err != nil {
			return err
		}
		if bp != "" {
			key = []byte(bp)
		}
		if multi {
			return batchOp(tx, b, true)
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
		WithCompleter(hintKey, hintPath)
	Cmd("mv", "Move (rename) a key").WithParams("src", "dst").WithHandler(handleMv).
		WithCompleter(hintKey, hintPath)
}
