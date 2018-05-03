package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

type Task struct {
	Description string
	Status      string
	CreatedAt   time.Time
	ClosedAt    time.Time
	tasklist    *TaskList
	next        *Task
	prev        *Task
}

type TaskList struct {
	first  *Task
	last   *Task
	length int
}

func (task *Task) Next() *Task {
	return task.next
}

func (task *Task) Prev() *Task {
	return task.prev
}

func (task *Task) checkIsElement() {
	if task.tasklist == nil {
		panic("Receiver is not a tasklist element")
	}
}

func (task *Task) InsertBefore(newTask *Task) *Task {
	task.checkIsElement()
	newTask.tasklist = task.tasklist
	newTask.next = task
	if task.prev == nil {
		newTask.prev = nil
		task.prev = newTask
		task.tasklist.first = newTask
	} else {
		newTask.prev = task.prev
		task.prev.next = newTask
		task.prev = newTask
	}
	task.tasklist.length++
	return newTask
}

func (task *Task) InsertAfter(newTask *Task) *Task {
	task.checkIsElement()
	newTask.tasklist = task.tasklist
	newTask.prev = task
	if task.next == nil {
		newTask.next = nil
		task.next = newTask
		task.tasklist.last = newTask
	} else {
		task.next.prev = newTask
		newTask.next = task.next
		task.next = newTask
	}
	task.tasklist.length++
	return newTask
}

func (task *Task) Delete() {
	task.checkIsElement()
	if task.prev == nil {
		task.tasklist.first = task.next
	} else {
		task.prev.next = task.next
	}
	if task.next == nil {
		task.tasklist.last = task.prev
	} else {
		task.next.prev = task.prev
	}
	task.tasklist.length--
	// Prevent memleaks
	task.tasklist = nil
	task.next = nil
	task.prev = nil
}

func (tasklist *TaskList) Clear() {
	tasklist.first = nil
	tasklist.last = nil
	tasklist.length = 0
}

func (tasklist *TaskList) First() *Task {
	return tasklist.first
}

func (tasklist *TaskList) Last() *Task {
	return tasklist.last
}

func (tasklist *TaskList) Length() int {
	return tasklist.length
}

func (tasklist *TaskList) Append(task *Task) *Task {
	// Don't know where task came from
	task.tasklist = tasklist
	task.next = nil
	if tasklist.first == nil {
		task.prev = nil
		tasklist.first = task
		tasklist.last = task
	} else {
		task.prev = tasklist.last
		tasklist.last.next = task
		tasklist.last = task
	}
	tasklist.length++
	return task
}

func (tasklist *TaskList) Load(path string) error {
	tasklist.Clear()

	var tasks []Task
	yml, err := ioutil.ReadFile(path)
	if err == nil {
		err = yaml.Unmarshal(yml, &tasks)
	} else if os.IsNotExist(err) {
		// It's ok create file
		err = nil
	}

	// NOTE Go uses a copy of the value instead of the value
	// itself within a range clause, So always use index
	// when iterating slices to make TaskList operations
	for i, _ := range tasks {
		tasklist.Append(&tasks[i])
	}
	return err
}

func (tasklist *TaskList) Save(path string) error {
	var tasks []Task
	for t := tasklist.First(); t != nil; t = t.Next() {
		tasks = append(tasks, *t)
	}

	yml, err := yaml.Marshal(&tasks)
	if err == nil {
		err = ioutil.WriteFile(path, yml, 0644)
	}

	return err
}
