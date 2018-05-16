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

func startTermbox() {
	err := termbox.Init()
	check(err)
	termbox.SetInputMode(termbox.InputEsc)
	termbox.HideCursor()
}

func getCell(x, y int) termbox.Cell {
	w, _ := termbox.Size()
	buf := termbox.CellBuffer()
	return buf[y*w+x]
}

func setCellColors(x, y int, fg, bg termbox.Attribute) {
	termbox.SetCell(x, y, getCell(x, y).Ch, fg, bg)
}

func (tv *TaskView) render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// TODO Optimization: Move block below to separate event
	w, h := termbox.Size()
	tv.w = int(w/2) - 2 // minus padding
	tv.h = h - 5        // minus title, help line, and margins
	tv.x = 1
	tv.y = 3

	var b strings.Builder

	fmt.Fprintf(&b, "%s Tasks(%d)", tv.filter, len(tv.Tasks))
	editbox.Label(1, 1, 0, 0, 0, b.String())

	var prefix string
	var fg, bg termbox.Attribute
	for i, t := range tv.Page() {
		if t.IsClosed() {
			prefix = "C "
		} else {
			prefix = "• "
		}
		if i == tv.CursorToPage() {
			fg, bg = 0|termbox.AttrReverse, 0|termbox.AttrReverse
		} else {
			fg, bg = 0, 0
		}
		editbox.Label(0+tv.x, i+tv.y, tv.maxTaskW+2, fg, bg, prefix+t.Description)
	}
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
			case ev.Key == termbox.KeyInsert || ev.Ch == 'i':
				tv.InsertTaskBefore()
			case ev.Ch == 'a':
				tv.InsertTaskAfter()
			case ev.Key == termbox.KeyEnter || ev.Ch == 'e':
				t := tv.SelectedTask()
				if t != nil {
					tv.EditTask()
				} else {
					tv.AppendTask()
				}
			case ev.Key == termbox.KeyDelete || ev.Ch == 'd':
				tv.DeleteTask()
			case ev.Ch == 'c':
				tv.CloseTask()
			case ev.Ch == 'o':
				tv.ReopenTask()
			case ev.Ch == 'C':
				tv.Filter("Closed")
			case ev.Ch == 'O':
				tv.Filter("Open")
			case ev.Ch == 'A':
				tv.Filter("All")
			case ev.Key == termbox.KeyEsc || ev.Ch == 'q':
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
		fmt.Println("Usage: tasks YAMLFILE")
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

	startTermbox()
	tv := NewTaskView(tasklist)
	tv.Filter("Open")
	tv.render()
	tv.mainLoop()

	termbox.Close()
	err = tasklist.Save(filename)
	check(err)
}
