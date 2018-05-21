package main

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
	"testing"
)

// ----------------------------------------------------------------------------
// Support
// ----------------------------------------------------------------------------

func tvFixture(size int) *TaskView {
	tl := &TaskList{}
	data := []string{
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
	for i := 0; i < size; i++ {
		tl.Append(Task{Description: data[i], Status: StatusOpen})
	}
	tv := NewTaskView(tl)
	tv.h = len(tv.view)
	return tv
}

// ----------------------------------------------------------------------------

func TestTVNew(t *testing.T) {
	tl := &TaskList{}
	tv := NewTaskView(tl)
	assert.Equal(t, tv.cursor, 0)
	i, task := tv.SelectedTask()
	assert.Equal(t, i, -1)
	assert.True(t, task == nil)
}

func TestTVAppend(t *testing.T) {
	tv := tvFixture(3)
	assert.Equal(t, tv.String(), heredoc.Doc(`
	> [ ] foo
	  [ ] bar
	  [ ] baz
	`))
}

func TestTVInsertAndFilter(t *testing.T) {
	tv := tvFixture(3)
	tv.h = 100
	tl := tv.tasklist
	tl.InsertBefore(1, Task{Description: "qux", Status: StatusClosed})
	tl.InsertAfter(2, Task{Description: "quux", Status: StatusClosed})
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

func TestTVScrollingAndPaging(t *testing.T) {
	tv := tvFixture(22)
	tv.h = 5
	assert.Equal(t, tv.cursor, 0)
	tv.CursorDown()
	tv.CursorDown()
	assert.Equal(t, tv.String(), heredoc.Doc(`
	  [ ] foo
	  [ ] bar
	> [ ] baz
	  [ ] qux
	  [ ] quux
	`))
	i, task := tv.SelectedTask()
	assert.Equal(t, i, 2)
	assert.Equal(t, task.Description, "baz")
	tv.CursorDown()
	tv.CursorDown()
	tv.CursorDown()
	tv.CursorDown()
	assert.Equal(t, tv.cursor, 6)
	assert.Equal(t, tv.String(), heredoc.Doc(`
	  [ ] baz
	  [ ] qux
	  [ ] quux
	  [ ] corge
	> [ ] grault
	`))
	i, task = tv.SelectedTask()
	assert.Equal(t, i, 6)
	assert.Equal(t, task.Description, "grault")
	tv.PageDown()
	tv.PageDown()
	assert.Equal(t, tv.String(), heredoc.Doc(`
	  [ ] plugh
	  [ ] xyzzy
	  [ ] thud
	  [ ] blep
	> [ ] blah
	`))
	_, task = tv.SelectedTask()
	assert.Equal(t, task.Description, "blah")
	// Go to end
	tv.PageDown()
	tv.PageDown()
	tv.PageDown()
	tv.PageDown()

	assert.Equal(t, tv.String(), heredoc.Doc(`
	  [ ] wibble
	  [ ] wobble
	  [ ] wubble
	  [ ] flob
	> [ ] toto
	`))
	_, task = tv.SelectedTask()
	assert.Equal(t, task.Description, "toto")
	tv.CursorUp()
	tv.CursorUp()
	assert.Equal(t, tv.String(), heredoc.Doc(`
	  [ ] wibble
	  [ ] wobble
	> [ ] wubble
	  [ ] flob
	  [ ] toto
	`))
	_, task = tv.SelectedTask()
	assert.Equal(t, task.Description, "wubble")
	tv.PageUp()
	tv.PageUp()
	assert.Equal(t, tv.String(), heredoc.Doc(`
	> [ ] xyzzy
	  [ ] thud
	  [ ] blep
	  [ ] blah
	  [ ] boop
	`))
	_, task = tv.SelectedTask()
	assert.Equal(t, task.Description, "xyzzy")
}

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
