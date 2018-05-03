package main

import (
	"flag"
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
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

func startTermbox() {
	err := termbox.Init()
	check(err)
	termbox.SetInputMode(termbox.InputEsc)
	termbox.HideCursor()
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
	view := NewTaskView(tasklist)

	w, _ := termbox.Size()
	for i, t := range view.Tasks {
		printStrAt(0, i, w, t.Description, 0, 0)
	}
	termbox.Flush()

	termbox.PollEvent()
	termbox.Close()

}
