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
	return strings.TrimSpace(s[len(CommentPrefix) : len(s)-len(CommentSuffix)])
}

func MakeComment(s string) string {
	return fmt.Sprintf("%s %s %s", CommentPrefix, s, CommentSuffix)
}
