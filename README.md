## TaskBox

Yet another task manager for terminal. In fact this is a hybrid between
simple editor and simple task manager.

## Why?

I surrender my search for perfect task manager. I practice in Go.

## Features

  * TUI
  * tasklists in [github friendly markdown](TODO.md) format
  * no database backend
  * no lib deps
  * filters
  * undo/redo
  * archive
  * autosave

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
