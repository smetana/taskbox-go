package main

import (
	"flag"
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
	"strconv"
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

func printStrAt(x, y, w int, s string, fg, bg termbox.Attribute) {
	// We cannot rely on range index because in UTF-8 string
	// it shows byte position, not rune position
	i := -1
	for _, r := range s {
		i++
		if i >= w {
			break
		}
		termbox.SetCell(x+i, y, r, fg, bg)
	}
	for i = i + 1; i < w; i++ {
		termbox.SetCell(x+i, y, ' ', fg, bg)
	}
}

func (tv *TaskView) render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// TODO Optimization: Move block below to separate event
	w, h := termbox.Size()
	tv.w = int(w/2) - 3 // " > * Description"
	tv.h = h - 5        // minus title, help line, and margins
	x0, y0 := 3, 3      // tasklist position

	printStrAt(1, 1, tv.w, "Tasks ("+strconv.Itoa(len(tv.Tasks))+")", 0, 0)
	for i, t := range tv.Page() {
		printStrAt(0+x0, i+y0, tv.w, "• "+t.Description, 0, 0)
	}
	// Cursor
	printStrAt(x0-2, y0+tv.CursorToPage(), 1, ">", 0, 0)
	termbox.Flush()
}

func (tv *TaskView) mainLoop() {
	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowDown:
				tv.CursorDown()
			case termbox.KeyArrowUp:
				tv.CursorUp()
			case termbox.KeyPgdn:
				tv.PageDown()
			case termbox.KeyPgup:
				tv.PageUp()
			default:
				if ev.Ch != 0 {
					switch ev.Ch {
					case 'q', 'Q':
						return // Quit
					default:
						// do nothing
					}
				}
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

	tv.render()
	tv.mainLoop()

	termbox.Close()
}
