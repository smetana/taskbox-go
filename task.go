package main

import (
	"fmt"
	"strings"
)

const TaskPrefix string = "- [ ] "

type Status rune

const (
	StatusAll    Status = 0
	StatusOpen          = ' '
	StatusClosed        = 'x'
)

var statusToString = map[Status]string{
	StatusAll:    "All",
	StatusOpen:   "Open",
	StatusClosed: "Closed",
}

func (s Status) String() string {
	return statusToString[s]
}

func StatusFromString(s string) Status {
	for k, v := range statusToString {
		if v == s {
			return k
		}
	}
	return StatusAll
}

type Task struct {
	Description string
	Status      Status
}

func (task *Task) String() string {
	return fmt.Sprintf("- [%c] %s", task.Status, task.Description)
}

func ParseTask(s string) Task {
	if lineTypeOf(s) != lineTask {
		panic("Not a Task: " + s)
	}
	t := Task{}
	i := strings.IndexRune(s, '[') // in bytes!
	t.Status = Status(s[i+1])
	if len(s) >= len(TaskPrefix) {
		t.Description = s[len(TaskPrefix):]
	} else {
		t.Description = ""
	}
	return t
}
