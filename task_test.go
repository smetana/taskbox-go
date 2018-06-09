package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTaskPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	ParseTask("foo")
}

func TestParseTask(t *testing.T) {
	task := ParseTask("- [ ] foo")
	assert.Equal(t, task, Task{
		Description: "foo",
		Status:      StatusOpen,
	})

	task = ParseTask("- [x] bar")
	assert.Equal(t, task, Task{
		Description: "bar",
		Status:      StatusClosed,
	})

	task = ParseTask("- [x] ")
	assert.Equal(t, task, Task{
		Description: "",
		Status:      StatusClosed,
	})

	task = ParseTask("- [x]")
	assert.Equal(t, task, Task{
		Description: "",
		Status:      StatusClosed,
	})
}
