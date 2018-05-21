package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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

type TaskList struct {
	Tasks []Task
}

func (tasklist *TaskList) String() string {
	var s strings.Builder
	for _, t := range tasklist.Tasks {
		fmt.Fprintln(&s, t.String())
	}
	return s.String()
}

func (tasklist *TaskList) Append(task Task) {
	tasklist.Tasks = append(tasklist.Tasks, task)
}

func (tasklist *TaskList) InsertBefore(i int, task Task) {
	tasklist.Tasks = append(tasklist.Tasks, Task{})
	copy(tasklist.Tasks[i+1:], tasklist.Tasks[i:])
	tasklist.Tasks[i] = task
}

func (tasklist *TaskList) InsertAfter(i int, task Task) {
	if i == len(tasklist.Tasks)-1 {
		tasklist.Append(task)
	} else {
		tasklist.InsertBefore(i+1, task)
	}
}

func (tasklist *TaskList) Delete(i int) Task {
	task := tasklist.Tasks[i]
	copy(tasklist.Tasks[i:], tasklist.Tasks[i+1:])
	tasklist.Tasks[len(tasklist.Tasks)-1] = Task{}
	tasklist.Tasks = tasklist.Tasks[:len(tasklist.Tasks)-1]
	return task
}

func (tasklist *TaskList) Load(path string) error {
	tasklist.Tasks = make([]Task, 0)

	f, err := os.Open(path)
	if os.IsNotExist(err) {
		// It's ok, Will create file
		return nil
	}
	// But be aware of other errors
	check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		t := Task{}
		t.Status = Status(s[1])
		t.Description = s[4:]
		tasklist.Append(t)
	}
	return scanner.Err()
}

func (tasklist *TaskList) Save(path string) error {
	f, err := os.Create(path)
	check(err)
	defer f.Close()
	return ioutil.WriteFile(path, []byte(tasklist.String()), 0644)
}
