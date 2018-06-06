package main

import (
	"flag"
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/smetana/editbox-go"
	"os"
	"strings"
	"time"
)

var autosaveInterval time.Duration

func help() {
	termbox.Clear(0, 0)
	var shortcuts = []struct{ key, desc string }{
		{"k,Up", "Cursor Up"},
		{"j,Down", "Cursor Down"},
		{"a,Enter", "Enter Edit Mode (in Task Mode)"},
		{"Esc", "Enter Task Mode (in Edit mode)"},
		{"Tab", "Insert Task Prefix \"[ ]\" (in Edit Mode)"},
		{"i,Ins", "Insert Task/Line"},
		{"d,Del", "Delete Task/Line"},
		{"Space", "Toggle Task"},
		{"<,Left", "Move Task/Line Up"},
		{">,Right", "Move Task/Line Down"},
		{"~,Tab", "Change Filter"},
		{"u", "Undo"},
		{"r", "Redo"},
		{"s,^S", "Save"},
		{"q,^Q", "Quit"},
	}
	for i, sc := range shortcuts {
		editbox.Label(1, i+1, 8, 0|termbox.AttrBold, 0,
			fmt.Sprintf("%8s", sc.key))
		editbox.Label(11, i+1, 0, 0, 0, sc.desc)
	}
	termbox.Flush()
	termbox.PollEvent()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

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

	if tb.editor != nil {
		tb.editor.Render()
	}

	tb.renderStatusLine()
	termbox.Flush()
}

func (tb *TaskBox) renderStatusLine() {
	w, h := termbox.Size()
	var s strings.Builder
	fmt.Fprintf(&s, " Mode:%s", tb.mode.String())
	fmt.Fprintf(&s, "; Filter:%s", tb.filter.String())
	if autosaveInterval > 0 {
		fmt.Fprintf(&s, "; Autosave:%.0fmin", autosaveInterval.Minutes())
	}
	if tb.modified {
		fmt.Fprintf(&s, "; Modified")
	} else {
		if tb.undo.Len > 0 {
			fmt.Fprintf(&s, "; Saved")
		}
	}
	editbox.Label(0, h-1, w, 0, 0, s.String())
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


func autosave(tb *TaskBox, d time.Duration) {
	for {
		<- time.After(d)
		if tb.modified {
			tb.Save(tb.path)
			tb.renderStatusLine()
			termbox.Flush()
		}
	}
}

func main() {
	flag.Usage = func() {
		fmt.Println("Usage:\n  taskbox [options] filename\n\nOptions:")
		flag.PrintDefaults()
		fmt.Println()
	}
	flagStatus := flag.String("status", "",
		"Filter by task status on start (All,Open,Closed)")
	flagAutosave := flag.Int("autosave", 0,
		"Autosave interval in minutes (0 = Disable)")
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	autosaveInterval = time.Duration(*flagAutosave) * time.Minute

	tb := &TaskBox{filter: StatusFromString(*flagStatus)}
	tb.undo = NewUndo(tb)

	filename := flag.Args()[0]
	tb.Load(filename)

	err := termbox.Init()
	check(err)
	termbox.SetInputMode(termbox.InputEsc)
	termbox.HideCursor()

	tb.render()

	if *flagAutosave > 0 {
		go autosave(tb, autosaveInterval)
	}

	tb.mainLoop()

	termbox.Close()
}
