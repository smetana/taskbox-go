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

func confirm(msg string) bool {
	w, h := termbox.Size()
	// Clear line
	editbox.Label(1, h-1, w, 0, 0, "")
	return editbox.Confirm(1, h-1, 0|termbox.AttrBold, 0, msg)
}

func (tv *TaskView) render() {
	termbox.Clear(0, 0)
	w, h := termbox.Size()
	tv.w = int(w/2) - 2 // minus margins
	tv.h = h - 4        // minus status and margins
	tv.x = 1
	tv.y = 1
	editbox.Text(tv.x, tv.y, 0, 0, 0, 0, tv.String())

	// status line
	var s strings.Builder
	var i int
	if len(tv.view) == 0 {
		i = 0
	} else {
		i = tv.cursor + 1
	}
	fmt.Fprintf(&s, "%d/%d ", i, len(tv.view))
	fmt.Fprintf(&s, "%s Tasks", tv.filter.String())
	if tv.tasklist.modified {
		fmt.Fprintf(&s, " - Modified")
	}
	editbox.Label(1, h-1, 0, 0, 0, s.String())

	termbox.Flush()
}

func (tv *TaskView) mainLoop() {
	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			switch {
			case ev.Key == termbox.KeyArrowDown ||
				ev.Ch == 'j':
				if tv.cursor == len(tv.view)-1 {
					tv.AppendTask()
				} else {
					tv.CursorDown()
				}
			case ev.Key == termbox.KeyArrowUp ||
				ev.Ch == 'k':
				tv.CursorUp()
			case ev.Key == termbox.KeyPgdn:
				tv.PageDown()
			case ev.Key == termbox.KeyPgup:
				tv.PageUp()
			case ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyEnter:
				if !menu(tv) {
					return
				}
			case ev.Key == termbox.KeyTab || ev.Ch == '~' || ev.Ch == '`':
				tv.NextFilter()
			case ev.Ch == '/':
				tv.AppendTask()
			case ev.Key == termbox.KeyInsert || ev.Ch == '+':
				tv.InsertTask()
			case ev.Ch == '\\':
				tv.EditTask()
			case ev.Key == termbox.KeyDelete || ev.Ch == '-':
				tv.DeleteTask()
			case ev.Key == termbox.KeySpace:
				tv.ToggleTask()
			case ev.Key == termbox.KeyArrowRight || ev.Ch == '>':
				tv.MoveTaskDown()
			case ev.Key == termbox.KeyArrowLeft || ev.Ch == '<':
				tv.MoveTaskUp()
			case ev.Key == termbox.KeyF1 ||
				ev.Ch == '?' ||
				ev.Ch == '/' ||
				ev.Ch == 'h' ||
				ev.Ch == 'H':
				help()
			case ev.Key == termbox.KeyCtrlQ ||
				ev.Key == termbox.KeyCtrlX ||
				ev.Ch == 'q' ||
				ev.Ch == 'Q':
				return // Quit
			default:
				// do nothing
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		tv.render()
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
	tasklist := &TaskList{}
	tasklist.Load(filename)

	err := termbox.Init()
	check(err)
	termbox.SetInputMode(termbox.InputEsc)
	termbox.HideCursor()

	tv := NewTaskView(tasklist)
	tv.Filter(StatusOpen)
	tv.render()
	tv.mainLoop()

	if tv.tasklist.modified && confirm("Save "+filename) {
		tasklist.Save(filename)
	}

	termbox.Close()
}
