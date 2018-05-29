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
      Esc  Menu
      k,↑  Cursor Up
      j,↓  Cursor Down
      \,/  Append Task
    Space  Toggle Status
    Enter  Edit
      <,←  Move Task Up
      >,→  Move Task Down
    +,Ins  Insert Task
    -,Del  Delete Task
    ~,Tab  Change Filter
     q,^Q  Quit
`)
	termbox.Flush()
	termbox.PollEvent()
}

func menu(tv *TaskView) bool {
	termbox.Clear(0, 0)
	menu := editbox.Select(
		2, 2, 20, 10,
		0, 0, 0|termbox.AttrReverse, 0|termbox.AttrReverse,
		[]string{
			"Continue",
			"",
			"Filter: Open Tasks",
			"Filter: Closed Tasks",
			"Filter: All Tasks",
			"",
			"Exit & Save",
			"Exit & Don't Save",
		},
	)
	ev := menu.WaitExit()
	if ev.Key == termbox.KeyEsc {
		return true
	}
	switch menu.Text() {
	case "Filter: Open Tasks":
		tv.Filter(StatusOpen)
	case "Filter: Closed Tasks":
		tv.Filter(StatusClosed)
	case "Filter: All Tasks":
		tv.Filter(StatusAll)
	case "Exit & Save":
		tv.tasklist.Save(tv.tasklist.path)
		return false
	case "Exit & Don't Save":
		tv.tasklist.modified = false
		return false
	}
	return true
}

func confirm(msg string) bool {
	termbox.Clear(0, 0)
	return editbox.Confirm(1, 1, 0, 0, msg)
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
				tv.CursorDown()
			case ev.Key == termbox.KeyArrowUp ||
				ev.Ch == 'k':
				tv.CursorUp()
			case ev.Key == termbox.KeyPgdn:
				tv.PageDown()
			case ev.Key == termbox.KeyPgup:
				tv.PageUp()
			case ev.Key == termbox.KeyEsc:
				if !menu(tv) {
					return
				}
			case ev.Key == termbox.KeyTab ||
				ev.Ch == '~' ||
				ev.Ch == '`':
				tv.NextFilter()
			case ev.Ch == 'n' ||
				ev.Ch == '/' ||
				ev.Ch == '\\':
				tv.AppendTask()
			case ev.Key == termbox.KeyInsert ||
				ev.Ch == '+':
				tv.InsertTask()
			case ev.Key == termbox.KeyEnter:
				tv.EditTask()
			case ev.Key == termbox.KeyDelete ||
				ev.Ch == '-':
				_, t := tv.SelectedTask()
				if t != nil && confirm("Delete \""+t.Description+"\"?") {
					tv.DeleteTask()
				}
			case ev.Key == termbox.KeySpace:
				tv.ToggleTask()
			case ev.Key == termbox.KeyArrowRight ||
				ev.Ch == '>':
				tv.MoveTaskDown()
			case ev.Key == termbox.KeyArrowLeft ||
				ev.Ch == '<':
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
