## TaskBox

Yet another task manager for terminal. In fact this is a hybrid between
simple editor and simple task manager.

## Why?

I surrender to find perfect (for me) task manager. I practice in Go.

## Features

  * TUI
  * uses plain text files in [github friendly](https://help.github.com/articles/about-task-lists/)
    markdown format (see [TODO.md](TODO.md))
  * does not require any libs
  * unlimited undo/redo
  * archive
  * autosave


## Screenshot

![Screenshot](screenshot.png)

## Installation

You will need to have [Go](https://golang.org) installed and configured
in your path.

```
git clone https://github.com/smetana/taskbox-go
cd taskbox-go
go get -d ./
go build -o taskbox
```

## Usage
```
./taskbox <filename>
```
