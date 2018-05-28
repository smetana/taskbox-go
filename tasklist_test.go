package main

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

// ----------------------------------------------------------------------------
// Support
// ----------------------------------------------------------------------------

func tasklistFixture() *TaskList {
	tl := &TaskList{}
	data := []string{"Foo", "Bar", "Baz"}
	for _, s := range data {
		tl.Append(Task{Description: s, Status: StatusOpen})
	}
	return tl
}

// ----------------------------------------------------------------------------

func TestTLAppend(t *testing.T) {
	tl := tasklistFixture()
	assert.Equal(t, tl.String(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[ ] Baz
	`))
}

func TestTLInsert(t *testing.T) {
	tl := tasklistFixture()
	tl.InsertAfter(1, Task{Description: "Qux", Status: StatusOpen})
	tl.InsertBefore(1, Task{Description: "Xyz", Status: StatusOpen})
	assert.Equal(t, tl.String(), heredoc.Doc(`
		[ ] Foo
		[ ] Xyz
		[ ] Bar
		[ ] Qux
		[ ] Baz
	`))
}

func TestTLSaveLoad(t *testing.T) {
	file, _ := ioutil.TempFile("", "test.yml")
	defer os.Remove(file.Name())

	tl1 := tasklistFixture()
	assert.Equal(t, tl1.String(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[ ] Baz
	`))
	tl1.Save(file.Name())

	tl2 := &TaskList{}
	tl2.Load(file.Name())

	assert.Equal(t, tl2.String(), heredoc.Doc(`
		[ ] Foo
		[ ] Bar
		[ ] Baz
	`))
}

func TestTLDelete(t *testing.T) {
	tl := tasklistFixture()
	task := tl.Delete(1)
	assert.Equal(t, tl.String(), heredoc.Doc(`
		[ ] Foo
		[ ] Baz
	`))
	assert.Equal(t, task.String(), "[ ] Bar")
}

func TestTLSwap(t *testing.T) {
	tl := tasklistFixture()
	tl.Swap(0, 2)
	assert.Equal(t, tl.String(), heredoc.Doc(`
		[ ] Baz
		[ ] Bar
		[ ] Foo
	`))
	tl.Swap(1, 0)
	assert.Equal(t, tl.String(), heredoc.Doc(`
		[ ] Bar
		[ ] Baz
		[ ] Foo
	`))
}
