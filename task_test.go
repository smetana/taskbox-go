package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTask(t *testing.T) {
	isTask, task := ParseTask("[ ] foo")
	assert.True(t, isTask)
	assert.Equal(t, task, Task{
		Description: "foo",
		Status:      StatusOpen,
	})

	isTask, task = ParseTask("[X] bar")
	assert.Equal(t, task, Task{
		Description: "bar",
		Status:      StatusClosed,
	})

	isTask, task = ParseTask("[X] ")
	assert.True(t, isTask)
	assert.Equal(t, task, Task{
		Description: "",
		Status:      StatusClosed,
	})

	isTask, task = ParseTask("[X]")
	assert.True(t, isTask)
	assert.Equal(t, task, Task{
		Description: "",
		Status:      StatusClosed,
	})

	isTask, task = ParseTask("[X")
	assert.False(t, isTask)

	isTask, task = ParseTask("")
	assert.False(t, isTask)

	isTask, task = ParseTask("  baz  ")
	assert.False(t, isTask)

	isTask, task = ParseTask("[@] ")
	assert.False(t, isTask)
}
