package main

import (
	"github.com/nsf/termbox-go"
	"github.com/smetana/editbox-go"
)

func (tb *TaskBox) EnterEditMode() {
	tb.mode = modeEdit
	index, _ := tb.SelectedLine()
	if index < 0 {
		tb.InsertLine(0, tb.TaskFilterPrefix())
		tb.calculate()
	}
	tb.AttachEditor()
	tb.editor.Render()
	tb.lastX, _ = tb.editor.GetCursor()
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
	if IsTask(tb.editor.Text()) {
		tb.editor.SetCursor(len(TaskPrefix), 0)
	} else {
		tb.editor.SetCursor(0, 0)
	}
}

func (tb *TaskBox) AddTaskPrefix() {
	pos, _ := tb.editor.GetCursor()
	if pos == 0 {
		if IsTask(tb.editor.Text()) {
			tb.editor.SetCursor(len(TaskPrefix), 0)
		} else {
			tb.editor.SetText(tb.TaskFilterPrefix())
		}
	}
	tb.editor.Render()
}

func (tb *TaskBox) EditBackspaceKey(ev termbox.Event) {
	pos, _ := tb.editor.GetCursor()
	if pos == 0 {
		i, _ := tb.SelectedLine()
		if i == 0 {
			return
		}
		s := tb.editor.Text()
		tb.DetachEditor()
		tb.DeleteLine(i)
		tb.calculate()
		tb.CursorUp()
		tb.render()
		tb.AttachEditor()
		tb.editor.SetText(s)
		tb.editor.SetCursor(len(tb.editor.Text())-len(s), 0)
	} else {
		ln := len(TaskPrefix)
		if pos == ln && IsTask(tb.editor.Text()) {
			for i := 0; i < ln; i++ {
				tb.editor.HandleEvent(ev)
			}
		} else {
			tb.editor.HandleEvent(ev)
		}
		tb.editor.Render()
	}
}

func (tb *TaskBox) InsertLineAndEdit() {
	i, s := tb.SelectedLine()
	var newLine string
	if IsTask(s) {
		newLine = tb.TaskFilterPrefix()
	} else {
		newLine = ""
	}
	tb.InsertLine(i, newLine)
	tb.calculate()
	tb.render()
	tb.EnterEditMode()
}

func (tb *TaskBox) EditMoveDown() {
	tb.DetachEditor()
	tb.CursorDown()
	tb.render()
	tb.AttachEditor()
	tb.editor.SetCursor(tb.lastX, 0)
}

func (tb *TaskBox) EditMoveUp() {
	tb.DetachEditor()
	tb.CursorUp()
	tb.render()
	tb.AttachEditor()
	tb.editor.SetCursor(tb.lastX, 0)
}

func (tb *TaskBox) EditMovePageDown() {
	tb.DetachEditor()
	tb.PageDown()
	tb.render()
	tb.AttachEditor()
	tb.editor.SetCursor(tb.lastX, 0)
}

func (tb *TaskBox) EditMovePageUp() {
	tb.DetachEditor()
	tb.PageUp()
	tb.render()
	tb.AttachEditor()
	tb.editor.SetCursor(tb.lastX, 0)
}

func (tb *TaskBox) HandleEditEvent(ev termbox.Event) {
	switch {
	case ev.Key == termbox.KeyEsc:
		tb.ExitEditMode()
	case ev.Key == termbox.KeyEnter:
		tb.EditEnterKey()
	case ev.Key == termbox.KeyTab:
		tb.AddTaskPrefix()
	case ev.Key == termbox.KeyBackspace || ev.Key == termbox.KeyBackspace2:
		tb.EditBackspaceKey(ev)
	case ev.Key == termbox.KeyArrowDown:
		tb.EditMoveDown()
	case ev.Key == termbox.KeyArrowUp:
		tb.EditMoveUp()
	case ev.Key == termbox.KeyPgdn:
		tb.EditMovePageDown()
	case ev.Key == termbox.KeyPgup:
		tb.EditMovePageUp()
	case ev.Key == termbox.KeyCtrlQ ||
		ev.Key == termbox.KeyCtrlX ||
		ev.Key == termbox.KeyCtrlC:
		tb.mode = modeExit
	default:
		tb.editor.HandleEvent(ev)
		// TODO Investigate why we need to render editor
		// to get correct cursor position
		tb.editor.Render()
		tb.lastX, _ = tb.editor.GetCursor()
	}
	if tb.editor != nil {
		tb.editor.Render()
		termbox.Flush()
	}
}
