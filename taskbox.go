package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/smetana/editbox-go"
	"strings"
)

type Mode int

const (
	modeTask Mode = iota
	modeEdit
	modeExit
)

func (m Mode) String() string {
	return map[Mode]string{
		modeTask: "Task",
		modeEdit: "Edit",
	}[m]
}

type TaskBox struct {
	mode     Mode
	Lines    []string
	path     string
	modified bool
	view     []int
	filter   Status
	x, y     int
	w, h     int
	cursor   int
	scroll   int
	editor   *editbox.Editbox
	lastX    int
}

func (tb *TaskBox) calculate() {
	tb.view = make([]int, 0)
	for i, line := range tb.Lines {
		isTask, t := ParseTask(line)
		if !isTask || t.Status == tb.filter || tb.filter == StatusAll {
			tb.view = append(tb.view, i)
		}
	}
	if len(tb.view) == 0 {
		tb.cursor = 0
	} else if tb.cursor >= len(tb.view) {
		tb.cursor = len(tb.view) - 1
	}
}

func (tb *TaskBox) Filter(s Status) {
	tb.cursor = 0
	tb.scroll = 0
	tb.filter = s
	tb.calculate()
}

func (tb *TaskBox) NextFilter() {
	filters := [3]Status{StatusOpen, StatusClosed, StatusAll}
	for i, f := range filters {
		if tb.filter == Status(f) {
			i++
			if i >= len(filters) {
				i = 0
			}
			tb.Filter(Status(filters[i]))
			return
		}
	}
}

func (tb TaskBox) TaskFilterPrefix() string {
	s := []rune(TaskPrefix)
	if tb.filter == StatusClosed {
		s[1] = StatusClosed
	}
	return string(s)
}

func (tb *TaskBox) String() string {
	if len(tb.view) == 0 {
		return "> No tasks. Press Enter to create one\n"
	}
	var to int
	if tb.scroll+tb.h > len(tb.view) {
		to = len(tb.view)
	} else {
		to = tb.scroll + tb.h
	}

	var s strings.Builder
	var cursor rune
	for i, index := range tb.view[tb.scroll:to] {
		if i == tb.CursorToPage() {
			cursor = '>'
		} else {
			cursor = ' '
		}
		fmt.Fprintf(&s, "%c %s\n", cursor, tb.Lines[index])
	}
	return s.String()
}

func (tb *TaskBox) scrollToCursor() {
	if tb.cursor-tb.scroll >= tb.h {
		tb.scroll = tb.cursor - tb.h + 1
	}
	if tb.cursor < tb.scroll {
		tb.scroll = tb.cursor
	}
}

func (tb *TaskBox) CursorDown() {
	if tb.cursor < len(tb.view)-1 {
		tb.cursor++
		tb.scrollToCursor()
	}
}

func (tb *TaskBox) CursorUp() {
	if tb.cursor > 0 {
		tb.cursor--
		tb.scrollToCursor()
	}
}

func (tb *TaskBox) PageDown() {
	tb.cursor = tb.cursor + tb.h - 1
	if tb.cursor >= len(tb.view) {
		tb.cursor = len(tb.view) - 1
	}
	tb.scrollToCursor()
}

func (tb *TaskBox) PageUp() {
	tb.cursor = tb.cursor - tb.h + 1
	if tb.cursor < 0 {
		tb.cursor = 0
	}
	tb.scrollToCursor()
}

func (tb *TaskBox) CursorToPage() int {
	return tb.cursor - tb.scroll
}

func (tb *TaskBox) CursorToY() int {
	return tb.y + tb.CursorToPage()
}

func (tb *TaskBox) SelectedLine() (int, string) {
	if len(tb.view) > 0 {
		index := tb.view[tb.cursor]
		return index, tb.Lines[index]
	} else {
		return -1, ""
	}
}

func (tb *TaskBox) HandleTaskEvent(ev termbox.Event) {
	switch {
	case ev.Key == termbox.KeyEnter || ev.Key == termbox.KeyEnd || ev.Ch == 'a':
		tb.EnterEditMode()
	case ev.Key == termbox.KeyInsert || ev.Ch == 'i':
		tb.InsertLineAndEdit()
	case ev.Key == termbox.KeyDelete || ev.Ch == 'd':
		tb.TaskDeleteKey()
	case ev.Key == termbox.KeyArrowDown || ev.Ch == 'j':
		tb.CursorDown()
	case ev.Key == termbox.KeyArrowUp || ev.Ch == 'k':
		tb.CursorUp()
	case ev.Key == termbox.KeyPgdn:
		tb.PageDown()
	case ev.Key == termbox.KeyPgup:
		tb.PageUp()
	case ev.Key == termbox.KeySpace:
		tb.ToggleTask()
	case ev.Key == termbox.KeyArrowLeft || ev.Ch == '<':
		tb.MoveLineUp()
	case ev.Key == termbox.KeyArrowRight || ev.Ch == '>':
		tb.MoveLineDown()
	case ev.Key == termbox.KeyTab || ev.Ch == '~' || ev.Ch == '`':
		tb.NextFilter()
	case ev.Key == termbox.KeyCtrlS || ev.Ch == 's':
		tb.Save(tb.path)
	case ev.Key == termbox.KeyF1 || ev.Ch == '?' || ev.Ch == 'h':
		help()
	case ev.Key == termbox.KeyCtrlQ ||
		ev.Key == termbox.KeyCtrlX ||
		ev.Key == termbox.KeyCtrlC ||
		ev.Ch == 'q':
		tb.mode = modeExit
	}
	tb.render()
}

func (tb *TaskBox) TaskDeleteKey() {
	i, _ := tb.SelectedLine()
	tb.DeleteLine(i)
	tb.calculate()
}

func (tb *TaskBox) ToggleTask() {
	i, s := tb.SelectedLine()
	isTask, task := ParseTask(s)
	if isTask {
		if task.Status == StatusOpen {
			task.Status = StatusClosed
		} else {
			task.Status = StatusOpen
		}
		tb.UpdateLine(i, task.String())
		tb.calculate()
	}
}

func (tb *TaskBox) MoveLineDown() {
	if tb.cursor >= len(tb.view)-1 {
		return
	}
	index1 := tb.view[tb.cursor]
	index2 := tb.view[tb.cursor+1]
	tb.SwapLines(index1, index2)
	tb.calculate()
	tb.CursorDown()
}

func (tb *TaskBox) MoveLineUp() {
	if tb.cursor <= 0 {
		return
	}
	index1 := tb.view[tb.cursor]
	index2 := tb.view[tb.cursor-1]
	tb.SwapLines(index1, index2)
	tb.calculate()
	tb.CursorUp()
}
