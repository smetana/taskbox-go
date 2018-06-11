package main

import (
	"github.com/nsf/termbox-go"
)

func (tb *TaskBox) EnterArchiveMode() {
	tb.mode = modeArchive
	tb.calculate()
}

func (tb *TaskBox) HandleArchiveEvent(ev termbox.Event) {
	switch {
	case ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlF:
		tb.mode = modeTask
	case ev.Key == termbox.KeyArrowDown || ev.Ch == 'j':
		tb.CursorDown()
	case ev.Key == termbox.KeyArrowUp || ev.Ch == 'k':
		tb.CursorUp()
	case ev.Key == termbox.KeyPgdn:
		tb.PageDown()
	case ev.Key == termbox.KeyPgup:
		tb.PageUp()
	case ev.Ch == 'z':
		tb.ToggleComment()
	case ev.Ch == 'u':
		tb.undo.Undo()
	case ev.Ch == 'r':
		tb.undo.Redo()
	case ev.Key == termbox.KeyCtrlS || ev.Ch == 's' || ev.Ch == 'w':
		tb.Save(tb.path)
	case ev.Ch == '?':
		help()
	case ev.Key == termbox.KeyCtrlQ ||
		ev.Key == termbox.KeyCtrlX ||
		ev.Key == termbox.KeyCtrlC ||
		ev.Ch == 'q':
		tb.mode = modeExit
	}
	tb.calculate()
	tb.render()
	termbox.Flush()
}
