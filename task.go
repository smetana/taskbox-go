package main

import (
	"fmt"
)

const TaskPrefix string = "[ ] "

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

func IsTask(s string) bool {
	return len(s) >= 3 &&
		s[0] == '[' &&
		s[2] == ']' &&
		(s[1] == StatusOpen || s[1] == StatusClosed)
}

func ParseTask(s string) (bool, Task) {
	t := Task{}
	if IsTask(s) {
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
