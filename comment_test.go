package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseCommentPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		} else {
			assert.Equal(t, "Not a comment: foo", r)
		}
	}()
	ParseComment("foo")
}

func TestParseComment(t *testing.T) {
	var pairs = []struct {
		s string
		c string
	}{
		{"<!-- Foo -->", "Foo"},
		{"<!--Foo -->", "Foo"},
		{"<!-- Foo-->", "Foo"},
		{"<!---->", ""},
	}
	for _, p := range pairs {
		c := ParseComment(p.s)
		if c != p.c {
			t.Errorf("got %s, want %s", c, p.s)
		}
	}
}

func TestMakeComment(t *testing.T) {
	c := MakeComment("Foo")
	assert.Equal(t, c, "<!-- Foo -->")
}
