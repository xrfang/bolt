package main

import (
    "fmt"
)

const (
	_G_HASH = "{_G_HASH}"
	_G_REVS = "{_G_REVS}"
    _BRANCH = "{_BRANCH}"	
	_G_OS = "{_G_OS}"
	_G_HW = "{_G_HW}"
)

func verinfo() string {
	ver := fmt.Sprintf("V%s.%s.%s.%s", _G_REVS, _G_HASH, _G_OS, _G_HW)
	if _BRANCH != "master" {
		ver += "/" + _BRANCH
	}	
	return ver
}
