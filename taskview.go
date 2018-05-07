package main

type TaskView struct {
	tasklist       *TaskList
	Tasks          []*Task
	filter         string
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
	}
	tv.scrollToCursor()
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

func (tv *TaskView) SelectedTask() *Task {
	return tv.Tasks[tv.cursor]
}
