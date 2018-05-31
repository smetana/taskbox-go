package main

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
	"testing"
)

// ----------------------------------------------------------------------------
// Support
// ----------------------------------------------------------------------------

var LinesFixture = []string{
	"foo",
	"bar",
	"baz",
	"qux",
	"quux",
	"corge",
	"grault",
	"garply",
	"waldo",
	"fred",
	"plugh",
	"xyzzy",
	"thud",
	"blep",
	"blah",
	"boop",
	"bloop",
	"wibble",
	"wobble",
	"wubble",
	"flob",
	"toto",
	"titi",
	"tata",
	"tutu",
}

func TaskBoxFixture(size int) *TaskBox {
	tb := &TaskBox{Lines: LinesFixture[0:size]}
	tb.calculate()
	tb.h = size
	return tb
}

// ----------------------------------------------------------------------------

func TestNewTaskBox(t *testing.T) {
	tb := NewTaskBox()
	assert.Equal(t, tb.cursor, 0)
	i, line := tb.SelectedLine()
	assert.Equal(t, i, -1)
	assert.True(t, line == "")
}

func TestTaskBoxString(t *testing.T) {
	tb := TaskBoxFixture(3)
	assert.Equal(t, tb.String(), heredoc.Doc(`
	> foo
	  bar
	  baz
	`))
}

func TestScrollingAndPaging(t *testing.T) {
	tb := TaskBoxFixture(22)
	tb.h = 5
	assert.Equal(t, tb.cursor, 0)
	tb.CursorDown()
	tb.CursorDown()
	assert.Equal(t, tb.String(), heredoc.Doc(`
	  foo
	  bar
	> baz
	  qux
	  quux
	`))
	i, line := tb.SelectedLine()
	assert.Equal(t, i, 2)
	assert.Equal(t, line, "baz")
	tb.CursorDown()
	tb.CursorDown()
	tb.CursorDown()
	tb.CursorDown()
	assert.Equal(t, tb.cursor, 6)
	assert.Equal(t, tb.String(), heredoc.Doc(`
	  baz
	  qux
	  quux
	  corge
	> grault
	`))
	i, line = tb.SelectedLine()
	assert.Equal(t, i, 6)
	assert.Equal(t, line, "grault")
	tb.PageDown()
	tb.PageDown()
	assert.Equal(t, tb.String(), heredoc.Doc(`
	  plugh
	  xyzzy
	  thud
	  blep
	> blah
	`))
	_, line = tb.SelectedLine()
	assert.Equal(t, line, "blah")
	// Go to end
	tb.PageDown()
	tb.PageDown()
	tb.PageDown()
	tb.PageDown()

	assert.Equal(t, tb.String(), heredoc.Doc(`
	  wibble
	  wobble
	  wubble
	  flob
	> toto
	`))
	_, line = tb.SelectedLine()
	assert.Equal(t, line, "toto")
	tb.CursorUp()
	tb.CursorUp()
	assert.Equal(t, tb.String(), heredoc.Doc(`
	  wibble
	  wobble
	> wubble
	  flob
	  toto
	`))
	_, line = tb.SelectedLine()
	assert.Equal(t, line, "wubble")
	tb.PageUp()
	tb.PageUp()
	assert.Equal(t, tb.String(), heredoc.Doc(`
	> xyzzy
	  thud
	  blep
	  blah
	  boop
	`))
	_, line = tb.SelectedLine()
	assert.Equal(t, line, "xyzzy")
}

/*
func TestTVMoveTaskDown(t *testing.T) {
	tv := tvFixture(3)
	tv.MoveTaskDown()
	assert.Equal(t, tv.String(), heredoc.Doc(`
	  [ ] bar
	> [ ] foo
	  [ ] baz
	`))
	tv.MoveTaskDown()
	assert.Equal(t, tv.String(), heredoc.Doc(`
	  [ ] bar
	  [ ] baz
	> [ ] foo
	`))
	tv.MoveTaskDown()
	assert.Equal(t, tv.String(), heredoc.Doc(`
	  [ ] bar
	  [ ] baz
	> [ ] foo
	`))
}

func TestTVInsertAndFilter(t *testing.T) {
	tv := tvFixture(3)
	tv.h = 100
	tl := tv.tasklist
	tl.Insert(1, Task{Description: "qux", Status: StatusClosed})
	tl.Insert(3, Task{Description: "quux", Status: StatusClosed})
	tv.calculate()
	assert.Equal(t, tv.String(), heredoc.Doc(`
	> [ ] foo
	  [X] qux
	  [ ] bar
	  [X] quux
	  [ ] baz
	`))
	tv.Filter(StatusOpen)
	assert.Equal(t, tv.String(), heredoc.Doc(`
	> [ ] foo
	  [ ] bar
	  [ ] baz
	`))
	tv.Filter(StatusClosed)
	assert.Equal(t, tv.String(), heredoc.Doc(`
	> [X] qux
	  [X] quux
	`))
	tv.Filter(StatusAll)
	assert.Equal(t, tv.String(), heredoc.Doc(`
	> [ ] foo
	  [X] qux
	  [ ] bar
	  [X] quux
	  [ ] baz
	`))
}

func TestTVMoveTaskUp(t *testing.T) {
	tv := tvFixture(3)
	tv.MoveTaskUp()
	tv.calculate()
	assert.Equal(t, tv.String(), heredoc.Doc(`
	> [ ] foo
	  [ ] bar
	  [ ] baz
	`))
	tv.CursorDown()
	tv.CursorDown()
	assert.Equal(t, tv.String(), heredoc.Doc(`
	  [ ] foo
	  [ ] bar
	> [ ] baz
	`))
	tv.MoveTaskUp()
	tv.calculate()
	assert.Equal(t, tv.String(), heredoc.Doc(`
	  [ ] foo
	> [ ] baz
	  [ ] bar
	`))
	tv.MoveTaskUp()
	tv.calculate()
	assert.Equal(t, tv.String(), heredoc.Doc(`
	> [ ] baz
	  [ ] foo
	  [ ] bar
	`))
	tv.MoveTaskUp()
	tv.calculate()
	assert.Equal(t, tv.String(), heredoc.Doc(`
	> [ ] baz
	  [ ] foo
	  [ ] bar
	`))
}

*/
