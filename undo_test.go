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
	tb.AppendLine("[X] Bar")
	tb.AppendLine("    Baz")

	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[X] Bar
		    Baz
	`))

	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[X] Bar
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
	tb.Lines = []string{"[ ] Foo", "[ ] Bar", "[X] Baz"}
	tb.InsertLine(2, "[X] Qux")
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
		[X] Qux
		[X] Baz
	`))
	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		Bar
		Foo
		[ ] Bar
		[X] Qux
		[X] Baz
	`))
	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		Foo
		Bar
		[ ] Bar
		[X] Qux
		[X] Baz
	`))
	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		FooBar
		[ ] Bar
		[X] Qux
		[X] Baz
	`))
	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[X] Qux
		[X] Baz
	`))
	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		## Xyz
		[ ] Bar
		[X] Qux
		[X] Baz
	`))
	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[X] Qux
		[X] Baz
	`))
	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[X] Baz
	`))
	tb.undo.Undo()
	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[X] Baz
	`))
}

func TestClearUndoOnLoad(t *testing.T) {
	file, _ := ioutil.TempFile("", "tasks.txt")
	defer os.Remove(file.Name())

	tb1 := TaskBoxWithUndo()
	tb1.Lines = []string{"[ ] Foo", "[ ] Bar", "[X] Baz"}
	assert.Equal(t, tb1.InnerString(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[X] Baz
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
		[X] Baz
	`))
}
