package main

import (
	"github.com/nsf/termbox-go"
	"github.com/smetana/editbox-go"
)

func (tb *TaskBox) EnterEditMode() {
	tb.mode = modeEdit
	tb.AttachEditor()
	tb.editor.Render()
}

func (tb *TaskBox) ExitEditMode() {
	tb.DetachEditor()
	termbox.HideCursor()
	tb.mode = modeTask
	tb.render()
	termbox.Flush()
}

// Attach editor at cursor
func (tb *TaskBox) AttachEditor() {
	_, s := tb.SelectedLine()
	tb.editor = editbox.Input(tb.x+2, tb.CursorToY(), tb.w-3, 0, 0)
	tb.editor.SetText(s)
}

func (tb *TaskBox) DetachEditor() {
	index, _ := tb.SelectedLine()
	tb.UpdateLine(index, tb.editor.Text())
	tb.editor = nil
}

func (tb *TaskBox) EditEnterKey() {
	pos, _ := tb.editor.GetCursor()
	tb.DetachEditor()
	i, _ := tb.SelectedLine()
	tb.SplitLine(i, pos)
	tb.calculate()
	tb.CursorDown()
	tb.render()
	tb.AttachEditor()
}

func (tb *TaskBox) HandleEditEvent(ev termbox.Event) {
	switch {
	case ev.Key == termbox.KeyEsc:
		tb.ExitEditMode()
	case ev.Key == termbox.KeyEnter:
		tb.EditEnterKey()
	case ev.Key == termbox.KeyArrowDown:
	case ev.Key == termbox.KeyArrowUp:
	case ev.Key == termbox.KeyPgdn:
	case ev.Key == termbox.KeyPgup:
	case ev.Key == termbox.KeyTab:
	case ev.Key == termbox.KeyCtrlQ ||
		ev.Key == termbox.KeyCtrlX ||
		ev.Key == termbox.KeyCtrlC:
		tb.mode = modeExit
	default:
		tb.editor.HandleEvent(ev)
	}
	if tb.editor != nil {
		tb.editor.Render()
		termbox.Flush()
	}
}

