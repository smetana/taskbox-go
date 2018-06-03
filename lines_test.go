package main

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestInnerString(t *testing.T) {
	tb := TaskBox{}
	assert.Equal(t, tb.InnerString(), "\n")

	tb = TaskBox{Lines: []string{"Foo", "[ ] Bar", "[X] Baz"}}
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		Foo
		[ ] Bar
		[X] Baz
	`))
}

func TestAppendLine(t *testing.T) {
	tb := TaskBox{}
	tb.AppendLine("[ ] Foo")
	tb.AppendLine("[X] Bar")
	tb.AppendLine("    Baz")

	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[X] Bar
		    Baz
	`))
}

func TestInsertLine(t *testing.T) {
	tb := TaskBox{Lines: []string{"[ ] Foo", "[ ] Bar", "[X] Baz"}}

	tb.InsertLine(2, "[X] Qux")
	tb.InsertLine(1, "## Xyz")
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		## Xyz
		[ ] Bar
		[X] Qux
		[X] Baz
	`))
}

func TestSaveLoad(t *testing.T) {
	file, _ := ioutil.TempFile("", "tasks.txt")
	defer os.Remove(file.Name())

	tb1 := TaskBox{Lines: []string{"[ ] Foo", "[ ] Bar", "[X] Baz"}}
	assert.Equal(t, tb1.InnerString(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[X] Baz
	`))
	tb1.Save(file.Name())
	assert.Equal(t, tb1.path, file.Name())

	tb2 := &TaskBox{}
	tb2.Load(file.Name())

	assert.Equal(t, tb2.InnerString(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[X] Baz
	`))
	assert.Equal(t, tb2.path, file.Name())
}

func TestDeleteLine(t *testing.T) {
	tb := TaskBox{Lines: []string{"[ ] Foo", "[ ] Bar", "[X] Baz"}}
	line := tb.DeleteLine(1)
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[X] Baz
	`))
	assert.Equal(t, line, "[ ] Bar")
}

func TestSwapLines(t *testing.T) {
	tb := TaskBox{Lines: []string{"[ ] Foo", "[ ] Bar", "[X] Baz"}}
	tb.SwapLines(0, 2)
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[X] Baz
		[ ] Bar
		[ ] Foo
	`))
	tb.SwapLines(1, 0)
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Bar
		[X] Baz
		[ ] Foo
	`))
}

func TestSplitLines(t *testing.T) {
	tb := TaskBox{Lines: []string{"[ ] FooBar", "[@] FooBaz"}}
	tb.SplitLine(0, 7)
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[@] FooBaz
	`))
	tb.SplitLine(2, 7)
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[@] Foo
		Baz
	`))
	tb.SplitLine(3, 3)
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[@] Foo
		Baz

	`))
}

func TestSplitFiltered(t *testing.T) {
	tb := TaskBox{Lines: []string{"[X] FooBar", "[ ] FooBaz"}}
	tb.Filter(StatusClosed)
	tb.SplitLine(0, 7)
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[X] Foo
		[X] Bar
		[ ] FooBaz
	`))
	tb.SplitLine(2, 7)
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		[X] Foo
		[X] Bar
		[ ] Foo
		[X] Baz
	`))
}
