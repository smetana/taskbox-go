package main

import (
	"github.com/nsf/termbox-go"
	"github.com/smetana/editbox-go"
	"time"
)

type TaskView struct {
	tasklist       *TaskList
	Tasks          []*Task
	filter         string
	x, y           int
	w, h           int
	cursor, scroll int
}

func NewTaskView(tasklist *TaskList) *TaskView {
	tv := new(TaskView)
	tv.tasklist = tasklist
	tv.Filter("All")
	return tv
}

func (tv *TaskView) calculate() {
	tv.Tasks = make([]*Task, 0)
	for t := tv.tasklist.First(); t != nil; t = t.Next() {
		if tv.filter == "All" || t.Status == tv.filter {
			tv.Tasks = append(tv.Tasks, t)
		}
	}
	if len(tv.Tasks) == 0 {
		tv.cursor = 0
	} else if tv.cursor >= len(tv.Tasks) {
		tv.cursor = len(tv.Tasks) - 1
	}
}

func (tv *TaskView) Filter(status string) {
	tv.cursor = 0
	tv.scroll = 0
	tv.filter = status
	tv.calculate()
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
	if tv.cursor < len(tv.Tasks)-1 {
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
	if tv.cursor >= len(tv.Tasks) {
		tv.cursor = len(tv.Tasks) - 1
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

func (tv *TaskView) Page() []*Task {
	var to int
	if tv.scroll+tv.h > len(tv.Tasks) {
		to = len(tv.Tasks)
	} else {
		to = tv.scroll + tv.h
	}
	return tv.Tasks[tv.scroll:to]
}

func (tv *TaskView) CursorToPage() int {
	return tv.cursor - tv.scroll
}

func (tv *TaskView) CursorToY() int {
	return tv.y + tv.CursorToPage()
}

func (tv *TaskView) SelectedTask() *Task {
	if len(tv.Tasks) > 0 {
		return tv.Tasks[tv.cursor]
	} else {
		return nil
	}
}

func (tv *TaskView) NewTask() *Task {
	task := &Task{}
	if tv.filter == "All" {
		task.Status = "Open"
	} else {
		task.Status = tv.filter
	}
	task.CreatedAt = time.Now()
	return task
}

func (tv *TaskView) DeleteTask() {
	task := tv.SelectedTask()
	if task != nil {
		task.Delete()
		tv.calculate()
	}
}

func (tv *TaskView) EditTask() (*Task, termbox.Event) {
	task := tv.SelectedTask()

	input := editbox.Input(tv.x+2, tv.CursorToY(), tv.w, 0, 0)
	input.SetText(task.Description)
	ev := input.WaitExit()

	if input.Text() == "" {
		tv.DeleteTask()
		task = nil
	} else {
		task.Description = input.Text()
	}
	tv.calculate()
	tv.render()
	termbox.HideCursor()
	return task, ev
}

func (tv *TaskView) InsertTaskBefore() {
	tv.SelectedTask().InsertBefore(tv.NewTask())
	tv.calculate()
	tv.render()
	tv.EditTask()
}

func (tv *TaskView) InsertTaskAfter() {
	if tv.cursor == len(tv.Tasks)-1 {
		tv.AppendTask()
	} else {
		tv.CursorDown()
		tv.InsertTaskBefore()
	}
}

func (tv *TaskView) AppendTask() {
	for {
		task := tv.NewTask()
		tv.tasklist.Append(task)
		tv.calculate()
		tv.cursor = len(tv.Tasks) - 1
		tv.scrollToCursor()
		tv.render()
		_, ev := tv.EditTask()
		if ev.Key == termbox.KeyEsc {
			return
		}
	}
}

func (tv *TaskView) CloseTask() {
	task := tv.SelectedTask()
	task.Status = "Closed"
	task.ClosedAt = time.Now()
}

func (tv *TaskView) ReopenTask() {
	task := tv.SelectedTask()
	task.Status = "Open"
	task.ClosedAt = time.Time{}
	task.ReopenAt = time.Now()
}
