package main

import (
	"fmt"
	"io/fs"
	"strings"

	"go.etcd.io/bbolt"
)

var (
	db  *bbolt.DB
	bkt []string
	err error
)

func openBoltDB(fp string) error {
	mode := fs.FileMode(0666)
	if readonly {
		mode = fs.FileMode(0444)
	}
	db, err = bbolt.Open(fp, mode, nil)
	if err == nil {
		bkt = []string{fp}
	}
	return err
}

func view(fn func(*bbolt.Tx) error) {
	if db == nil {
		fmt.Println("ERROR: no database (try 'open' or 'create')")
		return
	}
	err := db.View(fn)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
}

func changeDir(tx *bbolt.Tx) (*bbolt.Bucket, error) {
	if len(bkt) < 2 {
		return nil, nil //current dir is root
	}
	b := tx.Bucket([]byte(bkt[1]))
	if b == nil {
		return nil, fmt.Errorf("bucket '%s' does not exist", bkt[1])
	}
	for i, p := range bkt[2:] {
		b = b.Bucket([]byte(p))
		if b == nil {
			return nil, fmt.Errorf("bucket '%s' does not exist",
				strings.Join(bkt[1:2+i], "/"))
		}
	}
	return b, nil
}
