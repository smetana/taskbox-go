package main

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAchive(t *testing.T) {
	tb := TaskBoxWithUndo()
	tb.Lines = []string{
		"## Foo",
		"- [ ] Foo",
		"",
		"## Bar",
		"- [ ] Bar",
		"<!-- - [ ] Bar -->",
		"## Baz",
		"- [x] Baz",
		"<!-- - [ ] Baz -->",
	}
	tb.calculate()
	tb.h = 100
	assert.Equal(t, tb.String(), heredoc.Doc(`
		> ## Foo
		  - [ ] Foo
		  
		  ## Bar
		  - [ ] Bar
		  ## Baz
		  - [x] Baz
	`))
	tb.cursor = 2
	tb.ToggleComment()
	tb.cursor = 3
	tb.ToggleComment()
	tb.cursor = 2
	tb.ToggleComment()
	assert.Equal(t, tb.String(), heredoc.Doc(`
		  ## Foo
		  - [ ] Foo
		> ## Baz
		  - [x] Baz
	`))

	tb.mode = modeArchive
	tb.calculate()
	assert.Equal(t, tb.String(), heredoc.Doc(`
		  
		  ## Bar
		> - [ ] Bar
		  - [ ] Bar
		  - [ ] Baz
	`))
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		## Foo
		- [ ] Foo
		<!--  -->
		<!-- ## Bar -->
		<!-- - [ ] Bar -->
		<!-- - [ ] Bar -->
		## Baz
		- [x] Baz
		<!-- - [ ] Baz -->
	`))
	tb.undo.Undo()
	tb.undo.Undo()
	assert.Equal(t, tb.InnerString(), heredoc.Doc(`
		## Foo
		- [ ] Foo
		<!--  -->
		## Bar
		- [ ] Bar
		<!-- - [ ] Bar -->
		## Baz
		- [x] Baz
		<!-- - [ ] Baz -->
	`))
}
