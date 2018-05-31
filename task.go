package main

import (
	"fmt"
)

type Status rune

const (
	StatusAll    Status = '*'
	StatusOpen          = ' '
	StatusClosed        = 'X'
)

func (s Status) String() string {
	return map[Status]string{
		StatusAll:    "All",
		StatusOpen:   "Open",
		StatusClosed: "Closed",
	}[s]
}

type Task struct {
	Description string
	Status      Status
}

func (task *Task) String() string {
	return fmt.Sprintf("[%c] %s", task.Status, task.Description)
}

func ParseTask(s string) (bool, Task) {
	t := Task{}
	if len(s) >= 3 && s[0] == '[' && s[2] == ']' {
		t.Status = Status(s[1])
		if len(s) > 4 {
			t.Description = s[4:]
		} else {
			t.Description = ""
		}
		return true, t
	}
	return false, t
}
