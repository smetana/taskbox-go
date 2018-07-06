package main

import (
	"bufio"
	"os"
	"strings"
)

func (tb *TaskBox) InnerString() string {
	if tb.mode == modeEdit {
		index, oldL := tb.SelectedLine()
		newL := tb.editor.Text()
		if oldL != newL {
			tb.UpdateLine(index, newL)
		}
	}
	return strings.Join(tb.Lines, "\n") + "\n"
}

func (tb *TaskBox) AppendLine(line string) {
	tb.Lines = append(tb.Lines, line)
}

func (tb *TaskBox) InsertLine(i int, line string) {
	tb.Lines = append(tb.Lines, "")
	copy(tb.Lines[i+1:], tb.Lines[i:])
	tb.Lines[i] = line
}

func (tb *TaskBox) UpdateLine(i int, newL string) {
	oldL := tb.Lines[i]
	if oldL == newL {
		return
	}
	tb.Lines[i] = newL
}

func (tb *TaskBox) DeleteLine(i int) string {
	line := tb.Lines[i]
	copy(tb.Lines[i:], tb.Lines[i+1:])
	tb.Lines[len(tb.Lines)-1] = ""
	tb.Lines = tb.Lines[:len(tb.Lines)-1]
	return line
}

func (tb *TaskBox) SwapLines(i, j int) {
	tb.Lines[i], tb.Lines[j] = tb.Lines[j], tb.Lines[i]
}

// Split line and copy everything on right to new line below
// Return new line index
func (tb *TaskBox) SplitLine(i, pos int) int {
	runes := []rune(tb.Lines[i])
	right := string(runes[pos:])
	tb.UpdateLine(i, string(runes[0:pos]))
	if lineTypeOf(tb.Lines[i]) == lineTask {
		right = tb.TaskFilterPrefix() + right
	}
	i++
	tb.InsertLine(i, right)
	return i
}

func (tb *TaskBox) MakeLastLine(i int) {
	line := tb.DeleteLine(i)
	tb.AppendLine(line)
}

/*
We store archived tasks in multiline comment
On load we transform all these lines to separate
comment lines to make it easier to work with. e.g:

	foo
	<!--
	bar
	baz
	-->

will be transformed to

	foo
	<!-- bar -->
	<!-- baz -->

*/
func (tb *TaskBox) Load(path string) {
	tb.path = path
	tb.Lines = make([]string, 0)

	f, err := os.Open(path)
	if os.IsNotExist(err) {
		// It's ok, Will create file
		return
	}
	check(err)
	defer f.Close()

	hasUndo := (tb.undo != nil)
	tb.undo = nil // Disable Undo
	cmtBlck := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		switch {
		case lineTypeOf(s) == lineCommentOpen:
			cmtBlck = true
			continue
		case lineTypeOf(s) == lineCommentClose:
			cmtBlck = false
			continue
		case cmtBlck:
			s = MakeComment(s)
		default:
			// use line as is
		}
		tb.AppendLine(s)
	}
	check(scanner.Err())
	tb.calculate()
	tb.modified = false
	if hasUndo {
		tb.undo = NewUndo(tb) // New Clear Undo
	}
}

/*
On save we collect all commented out lines to one multiline comment
and save it at the end of the file

	foo
	<!-- bar -->
	baz
	<!-- qux -->

will be transformed to

	foo
	baz
	<!--
	bar
	qux
	-->
*/

func (tb *TaskBox) Save(path string) {
	var comments []string
	var err error

	tb.path = path
	f, err := os.Create(path)
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, s := range tb.Lines {
		if lineTypeOf(s) == lineComment {
			// collect comments to write the at the end
			comments = append(comments, ParseComment(s))
		} else {
			w.WriteString(s)
			w.WriteRune('\n')
		}
	}
	if len(comments) > 0 {
		w.WriteString("<!--\n")
		for _, s := range comments {
			w.WriteString(s)
			w.WriteRune('\n')
		}
		w.WriteString("-->\n")
	}
	w.Flush()
	tb.modified = false
}
