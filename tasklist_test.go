package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// ----------------------------------------------------------------------------
// Support
// ----------------------------------------------------------------------------

func (tasklist *TaskList) toString() string {
	var s strings.Builder
	fmt.Fprintf(&s, "\n")
	for t := tasklist.First(); t != nil; t = t.Next() {
		fmt.Fprintf(&s, "%s\n", t.Description)
	}
	return s.String()
}

func fixture() *TaskList {
	tl := &TaskList{}
	data := []string{"Foo", "Bar", "Baz"}
	for _, s := range data {
		tl.Append(&Task{Description: s, Status: "Open"})
	}
	return tl
}

// ----------------------------------------------------------------------------

func TestAppend(t *testing.T) {
	tl := fixture()
	assert.Equal(t, tl.Length(), 3)
	assert.Equal(t, tl.toString(), `
Foo
Bar
Baz
`)
}

func TestInsert(t *testing.T) {
	tl := fixture();
	task := tl.First().Next()
	task.InsertAfter(&Task{Description: "Qux", Status: "Open"})
	task.InsertBefore(&Task{Description: "Xyz", Status: "Open"})
	assert.Equal(t, tl.Length(), 5)
	assert.Equal(t, tl.toString(), `
Foo
Xyz
Bar
Qux
Baz
`)
}

func TestSaveLoad(t *testing.T) {
	file, _ := ioutil.TempFile("", "test.yml")
	defer os.Remove(file.Name())

	tl1 := fixture()
	assert.Equal(t, tl1.toString(), `
Foo
Bar
Baz
`)
	tl1.Save(file.Name())

	tl2 := &TaskList{}
	tl2.Load(file.Name())

	assert.Equal(t, tl2.toString(), `
Foo
Bar
Baz
`)
}

func TestDelete(t *testing.T) {
	tl := fixture();
	assert.Equal(t, tl.Length(), 3)
	task := tl.First().Next()
	task.InsertAfter(&Task{Description: "Qux", Status: "Open"})
	assert.Equal(t, tl.Length(), 4)
	task.Delete()
	assert.Equal(t, tl.Length(), 3)
	assert.Equal(t, tl.toString(), `
Foo
Qux
Baz
`)
	assert.Equal(t, task.Description, "Bar")
}
