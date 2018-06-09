package main

import (
	"fmt"
	"strings"
)

const (
	CommentPrefix string = "<!--"
	CommentSuffix string = "-->"
)

func ParseComment(s string) string {
	if lineTypeOf(s) != lineComment {
		panic(fmt.Sprintf("Not a comment: %s", s))
	}
	r := []rune(s)
	return strings.TrimSpace(string(r[4 : len(r)-3]))
}

func MakeComment(s string) string {
	return fmt.Sprintf("%s %s %s", CommentPrefix, s, CommentSuffix)
}
