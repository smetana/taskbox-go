package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/smetana/editbox-go"
	"strings"
)

type TaskView struct {
	tasklist       *TaskList
	view           []int
	filter         Status
	x, y           int
	w, h           int
	cursor, scroll int
}

func NewTaskView(tasklist *TaskList) *TaskView {
	tv := new(TaskView)
	tv.tasklist = tasklist
	tv.Filter(StatusAll)
	return tv
}

func (tv *TaskView) calculate() {
	tv.view = make([]int, 0)
	for i, t := range tv.tasklist.Tasks {
		if tv.filter == StatusAll || t.Status == tv.filter {
			tv.view = append(tv.view, i)
		}
	}
	if len(tv.view) == 0 {
		tv.cursor = 0
	} else if tv.cursor >= len(tv.view) {
		tv.cursor = len(tv.view) - 1
	}
}

func (tv *TaskView) Filter(s Status) {
	tv.cursor = 0
	tv.scroll = 0
	tv.filter = s
	tv.calculate()
}

func (tv *TaskView) NextFilter() {
	filters := [3]rune{' ', 'X', '*'}
	for i, f := range filters {
		if tv.filter == Status(f) {
			i++
			if i >= len(filters) {
				i = 0
			}
			tv.Filter(Status(filters[i]))
			return
		}
	}
}

func (tv *TaskView) String() string {
	var to int
	if tv.scroll+tv.h > len(tv.view) {
		to = len(tv.view)
	} else {
		to = tv.scroll + tv.h
	}

	var s strings.Builder
	var cursor rune
	var t Task
	for i, index := range tv.view[tv.scroll:to] {
		if i == tv.CursorToPage() {
			cursor = '>'
		} else {
			cursor = ' '
		}
		t = tv.tasklist.Tasks[index]
		fmt.Fprintf(&s, "%c %s\n", cursor, t.String())
	}
	return s.String()
}

func (tv *TaskView) scrollToCursor() {
	if tv.cursor-tv.scroll >= tv.h {
		tv.scroll = tv.cursor - tv.h + 1
	}
	if tv.cursor < tv.scroll {
		tv.scroll = tv.cursor
	}
}

func (tv *TaskView) CursorDown() {
	if tv.cursor < len(tv.view)-1 {
		tv.cursor++
		tv.scrollToCursor()
	} else {
		tv.AppendTask()
	}
}

func (tv *TaskView) CursorUp() {
	if tv.cursor > 0 {
		tv.cursor--
	}
	tv.scrollToCursor()
}

func (tv *TaskView) PageDown() {
	tv.cursor = tv.cursor + tv.h - 1
	if tv.cursor >= len(tv.view) {
		tv.cursor = len(tv.view) - 1
	}
	tv.scrollToCursor()
}

func (tv *TaskView) PageUp() {
	tv.cursor = tv.cursor - tv.h + 1
	if tv.cursor < 0 {
		tv.cursor = 0
	}
	tv.scrollToCursor()
}

func (tv *TaskView) CursorToPage() int {
	return tv.cursor - tv.scroll
}

func (tv *TaskView) CursorToY() int {
	return tv.y + tv.CursorToPage()
}

func (tv *TaskView) SelectedTask() (int, *Task) {
	if len(tv.view) > 0 {
		index := tv.view[tv.cursor]
		return index, &tv.tasklist.Tasks[index]
	} else {
		return -1, nil
	}
}

func (tv *TaskView) NewTask() Task {
	task := Task{}
	if tv.filter == StatusAll {
		task.Status = StatusOpen
	} else {
		task.Status = tv.filter
	}
	return task
}

func (tv *TaskView) DeleteTask() {
	index, task := tv.SelectedTask()
	if task != nil {
		tv.tasklist.Delete(index)
		tv.calculate()
	}
}

func (tv *TaskView) EditTask() (*Task, termbox.Event) {
	_, task := tv.SelectedTask()

	if task == nil {
		return tv.AppendTask()
	}

	// TODO Refactor to Update by TaskList
	oldDescription := task.Description
	input := editbox.Input(tv.x+6, tv.CursorToY(), tv.w-3, 0, 0)
	input.SetText(task.Description)
	ev := input.WaitExit()

	if input.Text() == "" {
		tv.DeleteTask()
		task = nil
		tv.tasklist.modified = true
	} else {
		task.Description = input.Text()
		if oldDescription != task.Description {
			tv.tasklist.modified = true
		}
	}

	tv.calculate()
	tv.render()
	termbox.HideCursor()

	return task, ev
}

func (tv *TaskView) InsertTask() (*Task, termbox.Event) {
	index, task := tv.SelectedTask()

	if task == nil {
		return tv.AppendTask()
	}

	tv.tasklist.Insert(index, tv.NewTask())
	tv.calculate()
	tv.render()
	return tv.EditTask()
}

func (tv *TaskView) AppendTask() (*Task, termbox.Event) {
	for {
		tv.tasklist.Append(tv.NewTask())
		tv.calculate()
		tv.cursor = len(tv.view) - 1
		tv.scrollToCursor()
		tv.render()
		task, ev := tv.EditTask()
		if ev.Key == termbox.KeyEsc {
			return task, ev
		}
	}
}

func (tv *TaskView) ToggleTask() {
	_, task := tv.SelectedTask()
	if task == nil {
		return
	}
	if task.Status == StatusOpen {
		task.Status = StatusClosed
	} else {
		task.Status = StatusOpen
	}
	tv.tasklist.modified = true
	tv.calculate()
}

func (tv *TaskView) MoveTaskDown() {
	if tv.cursor >= len(tv.view)-1 {
		return
	}
	index1 := tv.view[tv.cursor]
	index2 := tv.view[tv.cursor+1]
	tv.tasklist.Swap(index1, index2)
	tv.calculate()
	tv.CursorDown()
}

func (tv *TaskView) MoveTaskUp() {
	if tv.cursor <= 0 {
		return
	}
	index1 := tv.view[tv.cursor]
	index2 := tv.view[tv.cursor-1]
	tv.tasklist.Swap(index1, index2)
	tv.calculate()
	tv.CursorUp()
}
