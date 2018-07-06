package main

import (
	"reflect"
)

type UndoState struct {
	lines  []string
	cursor int
	filter Status
}

type Undo struct {
	tb         *TaskBox
	history    []UndoState
	stateIndex int
}

func NewUndo(tb *TaskBox) *Undo {
	u := &Undo{tb: tb, stateIndex: -1}
	u.PutState()
	return u
}

func (u *Undo) GetState() UndoState {
	lines := make([]string, len(u.tb.Lines))
	copy(lines, u.tb.Lines)
	return UndoState{
		cursor: u.tb.cursor,
		filter: u.tb.filter,
		lines:  lines,
	}
}

func (u *Undo) CurrentState() UndoState {
	return u.history[u.stateIndex]
}

func (u *Undo) RestoreState() {
	state := u.CurrentState()
	u.tb.cursor = state.cursor
	u.tb.filter = state.filter
	u.tb.Lines = make([]string, len(state.lines))
	copy(u.tb.Lines, state.lines)
	u.tb.modified = true
}

func (u *Undo) PutState() {
	if u.stateIndex >= 0 && reflect.DeepEqual(u.CurrentState().lines, u.tb.Lines) {
		return
	}
	u.tb.modified = true
	u.history = u.history[:u.stateIndex+1]
	u.history = append(u.history, u.GetState())
	u.stateIndex++
}

func (u *Undo) Undo() {
	if u.stateIndex == 0 {
		return
	}
	u.stateIndex--
	u.RestoreState()
}

func (u *Undo) Redo() {
	if u.stateIndex == len(u.history)-1 {
		return
	}
	u.stateIndex++
	u.RestoreState()
}
