package main

import (
	"encoding/hex"
	"strings"

	"github.com/c-bata/go-prompt"
)

func matchRune(str, pattern []rune) bool {
	for len(pattern) > 0 {
		switch pattern[0] {
		default:
			if len(str) == 0 || str[0] != pattern[0] {
				return false
			}
		case '?':
			if len(str) == 0 {
				return false
			}
		case '*':
			return matchRune(str, pattern[1:]) || (len(str) > 0 &&
				matchRune(str[1:], pattern))
		}
		str = str[1:]
		pattern = pattern[1:]
	}
	return len(str) == 0 && len(pattern) == 0
}

func wildcardMatch(pattern, subject string) bool {
	if pattern == "" {
		return subject == pattern
	}
	if pattern == "*" {
		return true
	}
	return matchRune([]rune(subject), []rune(pattern))
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

// only used in go-prompt completer
func hintMatch(key []byte, pattern string) *prompt.Suggest {
	target := string(key)
	if !isPrintable(target) {
		target = hex.EncodeToString(key)
	}
	if fuzzyMatch(pattern, target) {
		return &prompt.Suggest{Text: target}
	}
	return nil
}

// only used in arguments in commands like rm/cp/mv
func wildMatch(key []byte, pattern string) bool {
	target := string(key)
	if !isPrintable(target) {
		target = hex.EncodeToString(key)
	}
	return wildcardMatch(pattern, target)
}
