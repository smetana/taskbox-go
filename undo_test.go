package main

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TaskBoxWithUndo() *TaskBox {
	tb := &TaskBox{}
	tb.undo = NewUndo(tb)
	return tb
}

func TestUndoUndoAppend(t *testing.T) {
	tb := TaskBoxWithUndo()
	tb.AppendLine("[ ] Foo")
	tb.AppendLine("[x] Bar")
	tb.AppendLine("    Baz")

	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[x] Bar
		    Baz
	`))

	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[x] Bar
	`))

	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
	`))

	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), "\n")
}

func TestUndo(t *testing.T) {
	tb := TaskBoxWithUndo()
	tb.Lines = []string{"[ ] Foo", "[ ] Bar", "[x] Baz"}
	tb.InsertLine(2, "[x] Qux")
	tb.InsertLine(1, "## Xyz")
	tb.DeleteLine(1)
	tb.InsertLine(1, "FooBar")
	tb.SplitLine(1, 3)
	tb.SwapLines(1, 2)
	tb.DeleteLine(2)
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		Bar
		[ ] Bar
		[x] Qux
		[x] Baz
	`))
	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		Bar
		Foo
		[ ] Bar
		[x] Qux
		[x] Baz
	`))
	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		Foo
		Bar
		[ ] Bar
		[x] Qux
		[x] Baz
	`))
	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		FooBar
		[ ] Bar
		[x] Qux
		[x] Baz
	`))
	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[x] Qux
		[x] Baz
	`))
	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		## Xyz
		[ ] Bar
		[x] Qux
		[x] Baz
	`))
	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[x] Qux
		[x] Baz
	`))
	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[x] Baz
	`))
	tb.undo.Undo()
	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[x] Baz
	`))
}

func TestClearUndoOnLoad(t *testing.T) {
	file, _ := ioutil.TempFile("", "tasks.txt")
	defer os.Remove(file.Name())

	tb1 := TaskBoxWithUndo()
	tb1.Lines = []string{"[ ] Foo", "[ ] Bar", "[x] Baz"}
	assert.Equal(t, tb1.InnerString(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[x] Baz
	`))
	tb1.Save(file.Name())

	tb2 := TaskBoxWithUndo()
	tb2.AppendLine("1")
	tb2.AppendLine("2")
	tb2.Load(file.Name())

	tb2.undo.Undo()
	tb2.undo.Undo()
	assert.Equal(t, tb2.InnerString(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[x] Baz
	`))
}

func TestRedo(t *testing.T) {
	tb := TaskBoxWithUndo()
	tb.Lines = []string{"[ ] Foo", "[ ] Bar", "[x] Baz"}
	tb.InsertLine(2, "[x] Qux")
	tb.InsertLine(1, "## Xyz")
	tb.DeleteLine(1)
	tb.InsertLine(1, "FooBar")
	tb.SplitLine(1, 3)
	tb.SwapLines(1, 2)
	tb.DeleteLine(2)
	tb.undo.Undo()
	tb.undo.Undo()
	tb.undo.Undo()
	tb.undo.Undo()
	tb.undo.Undo()
	tb.undo.Undo()
	tb.undo.Undo()
	tb.undo.Undo() // extra
	tb.undo.Undo() // extra
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[x] Baz
	`))
	tb.undo.Redo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[x] Qux
		[x] Baz
	`))
	tb.undo.Redo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		## Xyz
		[ ] Bar
		[x] Qux
		[x] Baz
	`))
	tb.undo.Redo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[x] Qux
		[x] Baz
	`))
	// Clears Redo
	tb.InsertLine(1, "FooBar")
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		FooBar
		[ ] Bar
		[x] Qux
		[x] Baz
	`))
	tb.undo.Redo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		FooBar
		[ ] Bar
		[x] Qux
		[x] Baz
	`))
	tb.undo.Undo()
	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		## Xyz
		[ ] Bar
		[x] Qux
		[x] Baz
	`))
	tb.undo.Redo()
	tb.undo.Redo()
	tb.undo.Redo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		FooBar
		[ ] Bar
		[x] Qux
		[x] Baz
	`))
}
