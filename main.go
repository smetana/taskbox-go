package main

import (
	"flag"
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/smetana/editbox-go"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func help() {
	termbox.Clear(0, 0)
	editbox.Text(1, 0, 0, 0, 0, 0, `
           Outdated Help

Esc,Enter  Menu
      k,↑  Cursor Up
      j,↓  Cursor Down
        /  Append Task
        \  Edit Task
    +,Ins  Insert Task
    -,Del  Delete Task
    Space  Toggle Status
      <,←  Move Task Up
      >,→  Move Task Down
    ~,Tab  Change Filter
     q,^Q  Quit
`)
	termbox.Flush()
	termbox.PollEvent()
}

/*
type menuItem struct {
	text string
	fn   func() bool
}

func menu(tv *TaskView) bool {
	termbox.Clear(0, 0)
	_, task := tv.SelectedTask()
	editbox.Label(1, 1, 0, 0|termbox.AttrBold, 0, task.String())

	menu := []menuItem{
		{"Current Task: Toggle", func() bool {
			tv.ToggleTask()
			return true
		}},
		{"Current Task: Edit", func() bool {
			tv.render()
			tv.EditTask()
			return true
		}},
		{"Current Task: Delete", func() bool {
			tv.DeleteTask()
			return true
		}},
		{"Current Task: Move Down", func() bool {
			tv.MoveTaskDown()
			return true
		}},
		{"Current Task: Move Up", func() bool {
			tv.MoveTaskUp()
			return true
		}},
		{"", nil},
		{"New Task: Append", func() bool {
			tv.render()
			tv.AppendTask()
			return true
		}},
		{"New Task: Insert", func() bool {
			tv.render()
			tv.InsertTask()
			return true
		}},
		{"", nil},
		{"Tasklist: Show Open", func() bool {
			tv.Filter(StatusOpen)
			return true
		}},
		{"Tasklist: Show Closed", func() bool {
			tv.Filter(StatusClosed)
			return true
		}},
		{"Taskiist: Show All", func() bool {
			tv.Filter(StatusAll)
			return true
		}},
		{"", nil},
		{"Continue: Do Nothing", func() bool {
			return true
		}},
		{"Insert Text", func() bool {
			tv.InsertComment()
			return true
		}},
		{"Help", func() bool {
			help()
			return true
		}},
		{"", nil},
		{"Exit & Save", func() bool {
			tv.tasklist.Save(tv.tasklist.path)
			return false
		}},
		{"Exit & Don't Save", func() bool {
			tv.tasklist.modified = false
			return false
		}},
	}

	menuText := make([]string, 0)
	for _, item := range menu {
		menuText = append(menuText, item.text)
	}

	menuBox := editbox.Select(
		1, 3, 24, 20,
		0, 0, 0|termbox.AttrReverse, 0|termbox.AttrReverse,
		menuText,
	)
	ev := menuBox.WaitExit()
	if ev.Key == termbox.KeyEsc {
		return true
	}
	for _, item := range menu {
		if item.text == menuBox.Text() {
			return item.fn()
		}
	}
	panic("Should not happen")
}
*/

func confirm(msg string) (bool, termbox.Event) {
	w, h := termbox.Size()
	// Clear line
	editbox.Label(1, h-1, w, 0, 0, "")
	return editbox.Confirm(1, h-1, 0|termbox.AttrBold, 0, msg)
}

func (tb *TaskBox) render() {
	termbox.Clear(0, 0)
	w, h := termbox.Size()
	tb.w = int(w/2) - 2 // minus margins
	tb.h = h - 4        // minus status and margins
	tb.x = 1
	tb.y = 1
	editbox.Text(tb.x, tb.y, 0, 0, 0, 0, tb.String())

	// status line
	var s strings.Builder
	fmt.Fprintf(&s, "Mode:%s", tb.mode.String())
	fmt.Fprintf(&s, "; Filter:%s", tb.filter.String())
	if tb.modified {
		fmt.Fprintf(&s, "; Modified")
	}
	editbox.Label(1, h-1, 0, 0, 0, s.String())

	if tb.editor != nil {
		tb.editor.Render()
	}

	termbox.Flush()
}

func (tb *TaskBox) mainLoop() {
	for tb.mode != modeExit {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventError {
			panic(ev.Err)
		}
		if ev.Type == termbox.EventInterrupt {
			tb.mode = modeExit
		}
		switch tb.mode {
		case modeTask:
			tb.HandleTaskEvent(ev)
		case modeEdit:
			tb.HandleEditEvent(ev)
		}
		if tb.mode == modeExit && tb.modified {
			yes, ev := confirm("Save " + tb.path)
			if ev.Key == termbox.KeyEsc {
				tb.mode = modeTask
			} else if yes {
				tb.Save(tb.path)
			}
		}
	}
}

func main() {

	flagset := flag.NewFlagSet("tasks", flag.ExitOnError)
	flagset.Usage = func() {
		fmt.Println("Usage: tasks filename")
		flagset.PrintDefaults()
	}
	flagset.Parse(os.Args[1:])
	if len(flagset.Args()) == 0 {
		flagset.Usage()
		os.Exit(1)
	}
	filename := flagset.Args()[0]
	tb := &TaskBox{}
	tb.Load(filename)

	err := termbox.Init()
	check(err)
	termbox.SetInputMode(termbox.InputEsc)
	termbox.HideCursor()

	tb.Filter(StatusOpen)
	tb.render()
	tb.mainLoop()

	termbox.Close()
}
