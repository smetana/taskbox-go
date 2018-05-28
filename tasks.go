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

func getCell(x, y int) termbox.Cell {
	w, _ := termbox.Size()
	return termbox.CellBuffer()[y*w+x]
}

func setCellColors(x, y int, fg, bg termbox.Attribute) {
	cell := getCell(x, y)
	termbox.SetCell(x, y, cell.Ch, fg, bg)
}

func printHelp(x, y int, s string) {
	i := 0
	shortcut := false
	fg := termbox.ColorDefault
	for _, r := range s {
		if r == '_' {
			shortcut = !shortcut
		} else {
			if shortcut {
				fg = termbox.ColorDefault | termbox.AttrUnderline
			} else {
				fg = termbox.ColorDefault
			}
			termbox.SetCell(x+i, y, r, fg, 0)
			i++
		}
	}
}

func confirm(msg string) bool {
	termbox.Clear(0, 0)
	return editbox.Confirm(1, 1, 0, 0, msg)
}

func (tv *TaskView) render() {
	termbox.Clear(0, 0)

	// TODO Optimization: Move block below to separate event
	w, h := termbox.Size()
	tv.w = int(w/2) - 2 // minus padding
	tv.h = h - 5        // minus title, help line, and margins
	tv.x = 1
	tv.y = 4

	var title strings.Builder
	fmt.Fprintf(&title, "%s Tasks(%d)", tv.filter.String(), len(tv.view))
	if tv.modified {
		fmt.Fprintf(&title, " (modified)")
	}
	editbox.Label(1, 2, 0, 0, 0, title.String())

	editbox.Text(tv.x, tv.y, 0, 0, 0, 0, tv.String())

	printHelp(0, 0,
		"_m_enu  _n_ew  _i_nsert  _a_fter  _e_dit  "+
			"_d_elete  _t_oggle  _q_uit",
	)

	termbox.Flush()
}

func (tv *TaskView) mainLoop() {
	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			switch {
			case ev.Key == termbox.KeyArrowDown:
				tv.CursorDown()
			case ev.Key == termbox.KeyArrowUp:
				tv.CursorUp()
			case ev.Key == termbox.KeyPgdn:
				tv.PageDown()
			case ev.Key == termbox.KeyPgup:
				tv.PageUp()
			case ev.Key == termbox.KeyEsc ||
				ev.Ch == 'm' ||
				ev.Ch == 'M' ||
				ev.Ch == 'ь' ||
				ev.Ch == 'Ь':
				if !tv.ShowMenu() {
					return
				}
			case ev.Ch == 'n' ||
				ev.Ch == 'N' ||
				ev.Ch == 'т' ||
				ev.Ch == 'Т':
				tv.AppendTask()
			case ev.Key == termbox.KeyInsert ||
				ev.Ch == 'i' ||
				ev.Ch == 'I' ||
				ev.Ch == 'ш' ||
				ev.Ch == 'Ш':
				tv.InsertTaskBefore()
			case ev.Ch == 'a' ||
				ev.Ch == 'A' ||
				ev.Ch == 'ф' ||
				ev.Ch == 'Ф':
				tv.InsertTaskAfter()
			case ev.Key == termbox.KeyEnter ||
				ev.Ch == 'e' ||
				ev.Ch == 'E' ||
				ev.Ch == 'у' ||
				ev.Ch == 'У':
				tv.EditTask()
			case ev.Key == termbox.KeyDelete ||
				ev.Ch == 'd' ||
				ev.Ch == 'D' ||
				ev.Ch == 'в' ||
				ev.Ch == 'В':
				_, t := tv.SelectedTask()
				if t != nil && confirm("Delete \""+t.Description+"\"?") {
					tv.DeleteTask()
				}
			case ev.Key == termbox.KeySpace ||
				ev.Ch == 't' ||
				ev.Ch == 'T' ||
				ev.Ch == 'е' ||
				ev.Ch == 'Е':
				tv.ToggleTask()
			case ev.Key == termbox.KeyArrowRight:
				tv.MoveTaskDown()
			case ev.Key == termbox.KeyArrowLeft:
				tv.MoveTaskUp()
			case ev.Ch == 'q' ||
				ev.Ch == 'Q' ||
				ev.Ch == 'й' ||
				ev.Ch == 'Й':
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
	err := tasklist.Load(filename)
	check(err)

	err = termbox.Init()
	check(err)
	termbox.SetInputMode(termbox.InputEsc)
	termbox.HideCursor()

	tv := NewTaskView(tasklist)
	tv.Filter(StatusOpen)
	tv.render()
	tv.mainLoop()

	if tv.modified && confirm("Save "+filename) {
		err = tasklist.Save(filename)
		check(err)
	}

	termbox.Close()
}
