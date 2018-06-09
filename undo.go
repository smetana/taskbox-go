package main

type ChangeAction int

const (
	ActionInsert ChangeAction = iota
	ActionDelete
	ActionUpdate
	ActionSwap
)

type Change struct {
	Cursor     int
	Filter     Status
	Action     ChangeAction
	LineIndex  int
	OldLine    string
	NewLine    string
	LineIndex1 int // lines swap
	LineIndex2 int
}

func ChangeInsertLine(i int, s string) Change {
	return Change{
		Action:    ActionInsert,
		LineIndex: i,
		NewLine:   s,
	}
}

func ChangeUpdateLine(i int, oldL, newL string) Change {
	return Change{
		Action:    ActionUpdate,
		LineIndex: i,
		OldLine:   oldL,
		NewLine:   newL,
	}
}

func ChangeDeleteLine(i int, line string) Change {
	return Change{
		Action:    ActionDelete,
		LineIndex: i,
		OldLine:   line,
	}
}

func ChangeSwapLines(i, j int) Change {
	return Change{
		Action:     ActionSwap,
		LineIndex1: i,
		LineIndex2: j,
	}
}

type Undo struct {
	tb      *TaskBox
	Changes [][]Change
	Chain   *[]Change
	Len     int
}

func NewUndo(tb *TaskBox) *Undo {
	return &Undo{tb: tb}
}

func (u *Undo) StartChain() {
	chain := make([]Change, 0)
	u.Chain = &chain
}

func (u *Undo) PutChain() {
	u.Changes = append(u.Changes, *u.Chain)
	u.Chain = nil
	u.Len = len(u.Changes)
}

func (u *Undo) Put(c Change) {
	c.Cursor = u.tb.cursor
	c.Filter = u.tb.filter
	if u.Chain != nil {
		*u.Chain = append(*u.Chain, c)
	} else {
		u.Changes = append(u.Changes, []Change{c})
	}
	u.Len = len(u.Changes)
}

func (u *Undo) Undo() {
	if len(u.Changes) <= 0 {
		return
	}
	u.tb.undo = nil // Disable Undo
	chain := u.Changes[len(u.Changes)-1]
	// play chain backward
	for i := len(chain) - 1; i >= 0; i-- {
		change := chain[i]
		switch change.Action {
		case ActionInsert:
			u.tb.DeleteLine(change.LineIndex)
		case ActionDelete:
			u.tb.InsertLine(change.LineIndex, change.OldLine)
		case ActionUpdate:
			u.tb.UpdateLine(change.LineIndex, change.OldLine)
		case ActionSwap:
			u.tb.SwapLines(change.LineIndex1, change.LineIndex2)
		}
		u.tb.cursor = change.Cursor
		u.tb.filter = change.Filter

	}
	u.Changes = u.Changes[0 : len(u.Changes)-1]
	u.tb.undo = u
	u.tb.calculate()
}

func (u *Undo) Redo() {
	if u.Len == len(u.Changes) {
		return
	}
	u.Changes = u.Changes[0 : len(u.Changes)+1]
	u.tb.undo = nil // Disable Undo
	chain := u.Changes[len(u.Changes)-1]
	// play chain forward
	for _, change := range chain {
		switch change.Action {
		case ActionInsert:
			u.tb.InsertLine(change.LineIndex, change.NewLine)
		case ActionDelete:
			u.tb.DeleteLine(change.LineIndex)
		case ActionUpdate:
			u.tb.UpdateLine(change.LineIndex, change.NewLine)
		case ActionSwap:
			u.tb.SwapLines(change.LineIndex1, change.LineIndex2)
		}
		u.tb.cursor = change.Cursor
		u.tb.filter = change.Filter
	}
	u.tb.undo = u
	u.tb.calculate()
}
