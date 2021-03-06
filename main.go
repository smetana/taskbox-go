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
		{"k,Up", "cursor up"},
		{"j,Down", "cursor down"},
		{"Enter", "edit"},
		{"Esc", "stop edit"},
		{"Tab", "insert \"- [ ]\" (in Edit mode)"},
		{"i,Ins", "insert line"},
		{"d,Del", "delete line"},
		{"Space", "toggle status"},
		{"h,Left", "move line up"},
		{"l,Right", "move line down"},
		{"Ctrl+l", "move line to the bottom"},
		{"c", "insert copy of the line"},
		{"z", "archive line (unarchive line)"},
		{"f", "change filter"},
		{"Ctrl+f", "go to archive"},
		{"u", "undo"},
		{"r", "redo"},
		{"?", "help"},
		{"s,w", "save"},
		{"q", "quit"},
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
	tb.w = w - 2 // minus margins
	tb.h = h - 4 // minus status and margins
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
	if tb.mode != modeArchive {
		fmt.Fprintf(&s, "; Filter:%s", tb.filter.String())
	}
	if autosaveInterval > 0 {
		fmt.Fprintf(&s, "; Autosave:%.0fm", autosaveInterval.Minutes())
	}
	if tb.modified {
		fmt.Fprintf(&s, "; Modified")
	} else {
		if len(tb.undo.history) > 0 {
			fmt.Fprintf(&s, "; Saved  ")
		}
	}
	fmt.Fprintf(&s, "    %d:%d", tb.undo.stateIndex, len(tb.undo.history))
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
		case modeArchive:
			tb.HandleArchiveEvent(ev)
		}

		if !(ev.Ch == 'r' || ev.Ch == 'u') {
			tb.undo.PutState()
		}

		if tb.mode == modeExit && tb.modified {
			yes, ev := confirm("Save " + tb.path)
			if ev.Key == termbox.KeyEsc {
				tb.mode = modeTask
			} else if yes {
				tb.Save(tb.path)
			}
		}

		tb.calculate()
		tb.render()
	}
}

func autosave(tb *TaskBox, d time.Duration) {
	for {
		<-time.After(d)
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
