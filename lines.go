package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
)

func (tb *TaskBox) InnerString() string {
	return strings.Join(tb.Lines, "\n") + "\n"
}

func (tb *TaskBox) AppendLine(line string) {
	tb.Lines = append(tb.Lines, line)
	tb.modified = true
}

func (tb *TaskBox) InsertLine(i int, line string) {
	tb.Lines = append(tb.Lines, "")
	copy(tb.Lines[i+1:], tb.Lines[i:])
	tb.modified = true
	tb.Lines[i] = line
}

func (tb *TaskBox) UpdateLine(i int, newL string) {
	oldL := tb.Lines[i]
	if oldL != newL {
		tb.Lines[i] = newL
		tb.modified = true
	}
}

func (tb *TaskBox) DeleteLine(i int) string {
	line := tb.Lines[i]
	copy(tb.Lines[i:], tb.Lines[i+1:])
	tb.Lines[len(tb.Lines)-1] = ""
	tb.Lines = tb.Lines[:len(tb.Lines)-1]
	tb.modified = true
	return line
}

func (tb *TaskBox) SwapLines(i, j int) {
	tb.Lines[i], tb.Lines[j] = tb.Lines[j], tb.Lines[i]
	tb.modified = true
}

// Split line and copy everything on right to new line below
// Return new line index
func (tb *TaskBox) SplitLine(i, pos int) int {
	right := tb.Lines[i][pos:]
	tb.Lines[i] = tb.Lines[i][0:pos]
	i++
	tb.InsertLine(i, right)
	return i
}

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

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		tb.AppendLine(scanner.Text())
	}
	check(scanner.Err())
	tb.modified = false
}

func (tb *TaskBox) Save(path string) {
	tb.path = path
	f, err := os.Create(path)
	check(err)
	defer f.Close()
	err = ioutil.WriteFile(path, []byte(tb.InnerString()), 0644)
	check(err)
	tb.modified = false
}
