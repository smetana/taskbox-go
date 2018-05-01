package main

import (
	"container/list"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	//"strconv"
	"time"
)

type Task struct {
	Description string
	Status      string
	CreatedAt   time.Time
	ClosedAt    time.Time
}

type TaskList struct {
	list.List
}

func (tasklist *TaskList) Load(path string) error {
	tasklist.Init()

	var tasks []Task
	yml, err := ioutil.ReadFile(path)
	if err == nil {
		err = yaml.Unmarshal(yml, &tasks)
	} else if os.IsNotExist(err) {
		// It's ok create file
		err = nil
	}

	for _, t := range tasks {
		tasklist.PushBack(t)
	}
	return err
}

func (tasklist *TaskList) Save(path string) error {
	var tasks []Task
	for e := tasklist.Front(); e != nil; e = e.Next() {
		tasks = append(tasks, e.Value.(Task))
	}

	yml, err := yaml.Marshal(&tasks)
	if err == nil {
		err = ioutil.WriteFile(path, yml, 0644)
	}

	return err
}

func main() {
	tasklist := &TaskList{}
	tasklist.Load("test.yml")


	for e := tasklist.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value.(Task).Description)
	}
}
