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

	tb = TaskBox{Lines: []string{"Foo", "- [ ] Bar", "- [x] Baz"}}
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		Foo
		- [ ] Bar
		- [x] Baz
	`))
}

func TestAppendLine(t *testing.T) {
	tb := TaskBox{}
	tb.AppendLine("- [ ] Foo")
	tb.AppendLine("- [x] Bar")
	tb.AppendLine("      Baz")

	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		- [ ] Foo
		- [x] Bar
		      Baz
	`))
}

func TestInsertLine(t *testing.T) {
	tb := TaskBox{Lines: []string{"- [ ] Foo", "- [ ] Bar", "- [x] Baz"}}

	tb.InsertLine(2, "- [x] Qux")
	tb.InsertLine(1, "## Xyz")
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		- [ ] Foo
		## Xyz
		- [ ] Bar
		- [x] Qux
		- [x] Baz
	`))
}

func TestSaveLoad(t *testing.T) {
	file, _ := ioutil.TempFile("", "tasks.txt")
	defer os.Remove(file.Name())

	tb1 := TaskBox{Lines: []string{"- [ ] Foo", "- [ ] Bar", "- [x] Baz"}}
	assert.Equal(t, tb1.InnerString(), heredoc.Doc(`
		- [ ] Foo
		- [ ] Bar
		- [x] Baz
	`))
	tb1.Save(file.Name())
	assert.Equal(t, tb1.path, file.Name())

	tb2 := &TaskBox{}
	tb2.Load(file.Name())

	assert.Equal(t, tb2.InnerString(), heredoc.Doc(`
		- [ ] Foo
		- [ ] Bar
		- [x] Baz
	`))
	assert.Equal(t, tb2.path, file.Name())
}

func TestDeleteLine(t *testing.T) {
	tb := TaskBox{Lines: []string{"- [ ] Foo", "- [ ] Bar", "- [x] Baz"}}
	line := tb.DeleteLine(1)
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		- [ ] Foo
		- [x] Baz
	`))
	assert.Equal(t, line, "- [ ] Bar")
}

func TestSwapLines(t *testing.T) {
	tb := TaskBox{Lines: []string{"- [ ] Foo", "- [ ] Bar", "- [x] Baz"}}
	tb.SwapLines(0, 2)
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		- [x] Baz
		- [ ] Bar
		- [ ] Foo
	`))
	tb.SwapLines(1, 0)
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		- [ ] Bar
		- [x] Baz
		- [ ] Foo
	`))
}

func TestSplitLines(t *testing.T) {
	tb := TaskBox{Lines: []string{"- [ ] FooBar", "- [@] ФууБар"}}
	tb.SplitLine(0, 9)
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		- [ ] Foo
		- [ ] Bar
		- [@] ФууБар
	`))
	tb.SplitLine(2, 9)
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		- [ ] Foo
		- [ ] Bar
		- [@] Фуу
		Бар
	`))
	tb.SplitLine(3, 3)
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		- [ ] Foo
		- [ ] Bar
		- [@] Фуу
		Бар

	`))
}

func TestSplitFiltered(t *testing.T) {
	tb := TaskBox{Lines: []string{"- [x] FooBar", "- [ ] FooBaz"}}
	tb.Filter(StatusClosed)
	tb.SplitLine(0, 9)
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		- [x] Foo
		- [x] Bar
		- [ ] FooBaz
	`))
	tb.SplitLine(2, 9)
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		- [x] Foo
		- [x] Bar
		- [ ] Foo
		- [x] Baz
	`))
}

func TestMakeLastLine(t *testing.T) {
	tb := TaskBox{Lines: []string{
		"- [ ] Foo",
		"- [ ] Bar",
		"- [x] Baz",
		"- [ ] Qux",
	}}
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		- [ ] Foo
		- [ ] Bar
		- [x] Baz
		- [ ] Qux
	`))
	tb.MakeLastLine(1)
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		- [ ] Foo
		- [x] Baz
		- [ ] Qux
		- [ ] Bar
	`))
	tb.MakeLastLine(3)
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		- [ ] Foo
		- [x] Baz
		- [ ] Qux
		- [ ] Bar
	`))
}
