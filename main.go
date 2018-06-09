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
		{"Enter", "Enter Edit Mode"},
		{"Esc", "Back to Task Mode"},
		{"Tab", "Insert Task Prefix \"[ ]\" (in Edit Mode)"},
		{"i,Ins", "Insert Line"},
		{"d,Del", "Delete Line"},
		{"Space", "Toggle Status"},
		{"h,Left", "Move Line Up"},
		{"l,Right", "Move Line Down"},
		{"Ctrl+l", "Move Line to the Bottom"},
		{"z", "Archive/Unarchive Line"},
		{"f", "Change Filter"},
		{"Ctrl+f", "Show Archive"},
		{"u", "Undo"},
		{"r", "Redo"},
		{"?", "Help"},
		{"s,w", "Save"},
		{"q", "Quit"},
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
		case modeArchive:
			tb.HandleArchiveEvent(ev)
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
