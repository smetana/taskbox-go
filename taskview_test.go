package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
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
		tl.Append(&Task{Description: data[i], Status: "Open"})
	}
	tv := NewTaskView(tl)
	tv.h = len(tv.Tasks)
	return tv
}

func (tv *TaskView) toString() string {
	var s strings.Builder
	fmt.Fprintf(&s, "\n")
	for i, t := range tv.Tasks {
		fmt.Fprintf(&s, "%2d %-6s %s\n", i, t.Status, t.Description)
	}
	return s.String()
}

func (tv *TaskView) pageToString() string {
	var s strings.Builder
	var cursor string

	fmt.Fprintf(&s, "\n")

	for i, t := range tv.Page() {
		if i == tv.CursorToPage() {
			cursor = ">"
		} else {
			cursor = " "
		}
		fmt.Fprintf(&s, "%s %2d %-6s %s\n", cursor, tv.scroll+i, t.Status, t.Description)
	}

	return s.String()
}

// ----------------------------------------------------------------------------

func TestNew(t *testing.T) {
	tv := tvFixture(3)
	assert.Equal(t, tv.toString(), `
 0 Open   foo
 1 Open   bar
 2 Open   baz
`)
}

func TestInsertAndFilter(t *testing.T) {
	tv := tvFixture(3)
	task := tv.Tasks[1]
	task.InsertBefore(&Task{Description: "qux", Status: "Closed"})
	task.InsertAfter(&Task{Description: "quux", Status: "Closed"})
	tv.calculate()
	assert.Equal(t, tv.toString(), `
 0 Open   foo
 1 Closed qux
 2 Open   bar
 3 Closed quux
 4 Open   baz
`)
	tv.Filter("Open")
	assert.Equal(t, tv.toString(), `
 0 Open   foo
 1 Open   bar
 2 Open   baz
`)
	tv.Filter("Closed")
	assert.Equal(t, tv.toString(), `
 0 Closed qux
 1 Closed quux
`)
}

func TestScrollingAndPaging(t *testing.T) {
	tv := tvFixture(22)
	tv.h = 5
	assert.Equal(t, tv.cursor, 0)
	tv.CursorDown()
	tv.CursorDown()
	assert.Equal(t, tv.pageToString(), `
   0 Open   foo
   1 Open   bar
>  2 Open   baz
   3 Open   qux
   4 Open   quux
`)
	assert.Equal(t, tv.SelectedTask().Description, "baz")
	tv.CursorDown()
	tv.CursorDown()
	tv.CursorDown()
	tv.CursorDown()
	assert.Equal(t, tv.cursor, 6)
	assert.Equal(t, tv.pageToString(), `
   2 Open   baz
   3 Open   qux
   4 Open   quux
   5 Open   corge
>  6 Open   grault
`)
	assert.Equal(t, tv.SelectedTask().Description, "grault")
	tv.PageDown()
	tv.PageDown()
	assert.Equal(t, tv.pageToString(), `
  10 Open   plugh
  11 Open   xyzzy
  12 Open   thud
  13 Open   blep
> 14 Open   blah
`)
	assert.Equal(t, tv.SelectedTask().Description, "blah")
	// Go to end
	tv.PageDown()
	tv.PageDown()
	tv.PageDown()
	tv.PageDown()

	assert.Equal(t, tv.pageToString(), `
  17 Open   wibble
  18 Open   wobble
  19 Open   wubble
  20 Open   flob
> 21 Open   toto
`)
	assert.Equal(t, tv.SelectedTask().Description, "toto")
	tv.CursorUp()
	tv.CursorUp()
	assert.Equal(t, tv.pageToString(), `
  17 Open   wibble
  18 Open   wobble
> 19 Open   wubble
  20 Open   flob
  21 Open   toto
`)
	assert.Equal(t, tv.SelectedTask().Description, "wubble")
	tv.PageUp()
	tv.PageUp()
	assert.Equal(t, tv.pageToString(), `
> 11 Open   xyzzy
  12 Open   thud
  13 Open   blep
  14 Open   blah
  15 Open   boop
`)
	assert.Equal(t, tv.SelectedTask().Description, "xyzzy")
}
