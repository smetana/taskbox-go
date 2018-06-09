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
	lines := make([]string, size)
	copy(lines, LinesFixture[0:size])
	tb := &TaskBox{Lines: lines}
	tb.calculate()
	tb.h = size
	return tb
}

func TaskBoxWithUndo() *TaskBox {
	tb := &TaskBox{}
	tb.undo = NewUndo(tb)
	return tb
}

// ----------------------------------------------------------------------------

func TestLineTypeOf(t *testing.T) {
	var pairs = []struct {
		s string
		t lineType
	}{
		{"- [x", lineNormal},
		{"  baz  ", lineNormal},
		{"", lineNormal},
		{"- [ ] foo", lineTask},
		{"- [x] foo", lineTask},
		{"- [x] ", lineTask},
		{"- [x]", lineTask},
		{"- [@] foo", lineNormal},
		{"<!-- Foo -->", lineComment},
		{"<!-- Foo-->", lineComment},
		{"<!--Foo -->", lineComment},
		{"<!---->", lineComment},
		{"<!--->", lineNormal},
	}
	for _, p := range pairs {
		lt := lineTypeOf(p.s)
		if lt != p.t {
			t.Errorf("wrong type for %s got %d, want %d", p.s, lt, p.t)
		}
	}
}

func TestNewTaskBox(t *testing.T) {
	tb := &TaskBox{}
	assert.Equal(t, tb.cursor, 0)
	i, line := tb.SelectedLine()
	assert.Equal(t, i, -1)
	assert.True(t, line == "")
	assert.Equal(t, StatusAll, tb.filter)
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

func TestMoveLineDown(t *testing.T) {
	tb := TaskBoxFixture(3)
	tb.MoveLineDown()
	assert.Equal(t, tb.String(), heredoc.Doc(`
	  bar
	> foo
	  baz
	`))
	tb.MoveLineDown()
	assert.Equal(t, tb.String(), heredoc.Doc(`
	  bar
	  baz
	> foo
	`))
	tb.MoveLineDown()
	assert.Equal(t, tb.String(), heredoc.Doc(`
	  bar
	  baz
	> foo
	`))
}

func TestMoveLineUp(t *testing.T) {
	tb := TaskBoxFixture(3)
	assert.Equal(t, tb.cursor, 0)
	assert.Equal(t, tb.String(), heredoc.Doc(`
	> foo
	  bar
	  baz
	`))
	tb.MoveLineUp()
	assert.Equal(t, tb.String(), heredoc.Doc(`
	> foo
	  bar
	  baz
	`))
	tb.CursorDown()
	tb.CursorDown()
	assert.Equal(t, tb.String(), heredoc.Doc(`
	  foo
	  bar
	> baz
	`))
	tb.MoveLineUp()
	assert.Equal(t, tb.String(), heredoc.Doc(`
	  foo
	> baz
	  bar
	`))
	tb.MoveLineUp()
	assert.Equal(t, tb.String(), heredoc.Doc(`
	> baz
	  foo
	  bar
	`))
	tb.MoveLineUp()
	assert.Equal(t, tb.String(), heredoc.Doc(`
	> baz
	  foo
	  bar
	`))
}

func TestToggle(t *testing.T) {
	tb := &TaskBox{Lines: []string{
		"Foo",
		"- [ ] Bar",
		"- [x] Baz"}}
	tb.Filter(StatusAll)
	tb.calculate()
	tb.h = 3
	assert.Equal(t, tb.String(), heredoc.Doc(`
		> Foo
		  - [ ] Bar
		  - [x] Baz
	`))

	tb.CursorDown()
	tb.ToggleTask()
	assert.Equal(t, tb.String(), heredoc.Doc(`
		  Foo
		> - [x] Bar
		  - [x] Baz
	`))

	tb.CursorDown()
	tb.ToggleTask()
	assert.Equal(t, tb.String(), heredoc.Doc(`
		  Foo
		  - [x] Bar
		> - [ ] Baz
	`))
}

func TestToggleAndFilterOut(t *testing.T) {
	tb := &TaskBox{Lines: []string{
		"- [ ] Foo",
		"- [ ] Bar",
		"- [ ] Baz"}}
	tb.Filter(StatusOpen)
	tb.calculate()
	tb.h = 3
	assert.Equal(t, tb.String(), heredoc.Doc(`
		> - [ ] Foo
		  - [ ] Bar
		  - [ ] Baz
	`))

	tb.CursorDown()
	tb.ToggleTask()
	assert.Equal(t, tb.String(), heredoc.Doc(`
		  - [ ] Foo
		> - [ ] Baz
	`))

	tb.CursorDown()
	tb.ToggleTask()
	assert.Equal(t, tb.String(), heredoc.Doc(`
		> - [ ] Foo
	`))

	tb.Filter(StatusClosed)
	tb.calculate()
	tb.h = 3
	assert.Equal(t, tb.String(), heredoc.Doc(`
		> - [x] Bar
		  - [x] Baz
	`))

	tb.CursorDown()
	tb.ToggleTask()
	assert.Equal(t, tb.String(), heredoc.Doc(`
		> - [x] Bar
	`))

	tb.ToggleTask()
	assert.Equal(t, tb.String(), heredoc.Doc(`
		> No tasks. Press Enter to create one
	`))

}
