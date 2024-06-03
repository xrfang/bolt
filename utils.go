package main

import (
	"encoding/hex"
	"strings"
	"unicode"

	"github.com/c-bata/go-prompt"
)

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
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
