package main

import (
	"fmt"
	"strings"
	"time"

	"go.etcd.io/bbolt"
)

var (
	db  *bbolt.DB
	bkt []string
	err error
)

func openBoltDB(fp string) error {
	opts := bbolt.Options{
		Timeout:  time.Second,
		ReadOnly: readonly,
	}
	db, err = bbolt.Open(fp, 0666, &opts)
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

func hint(fn func(*bbolt.Tx) error) {
	if db == nil {
		return
	}
	db.View(fn)
}

func update(fn func(*bbolt.Tx) error) {
	if db == nil {
		fmt.Println("ERROR: no database (try 'open' or 'create')")
		return
	}
	err := db.Update(fn)
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
				strings.Join(bkt[1:3+i], "/"))
		}
	}
	return b, nil
}
