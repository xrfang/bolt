package main

import (
	"go.etcd.io/bbolt"
)

func init() {
	Cmd("cd", "Change current dir (open a bucket)").WithParams("dir").WithHandler(
		func(c *command) {
			dir := c.Arg("dir")
			switch dir {
			case "":
				bkt = bkt[:1]
			case "..":
				if l := len(bkt); l > 1 {
					bkt = bkt[:l-1]
				}
			default:
				view(func(tx *bbolt.Tx) error {
					bkt = append(bkt, dir)
					_, err = changeDir(tx)
					if err != nil {
						bkt = bkt[:len(bkt)-1]
					}
					return err
				})
			}
		},
	)
}
