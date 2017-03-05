package main

import (
	"github.com/nsf/termbox-go"
	"github.com/briansteffens/escapebox"
)

const (
	CommandMode = 0
	InsertMode  = 1
)

type Editbox struct {
	Bounds     Rect
	lines      []string
	cursorLine int
	cursorChar int
	scroll     int
	focus      bool
	mode       int
}

func splitRows(line string, textWidth int) []string {
	rows := len(line) / textWidth + 1
	ret := make([]string, rows)

	for i := 0; i < rows; i++ {
		start := i * textWidth
		stop := min((i + 1) * textWidth, len(line))
		ret[i] = line[start:stop]
	}

	return ret
}

func (e *Editbox) Render() {
	textWidth := e.Bounds.Width - 2
	textHeight := e.Bounds.Height - 3

	// Generate virtual lines and map the cursor to them.
	virtualLines := make([]string, 0)

	cursorRow := 0
	cursorCol := 0

	for lineIndex, line := range e.lines {
		virtualLineCount := len(line) / textWidth + 1

		if e.cursorLine == lineIndex {
			cursorRow = len(virtualLines) +
				    e.cursorChar / textWidth
			cursorCol = e.cursorChar % textWidth
		}

		for i := 0; i < virtualLineCount; i++ {
			start := i * textWidth
			stop := min(len(line), (i + 1) * textWidth)
			virtualLines = append(virtualLines, line[start:stop])
		}
	}

	if cursorRow < e.scroll {
		e.scroll = cursorRow
	}

	if cursorRow >= e.scroll + textHeight {
		e.scroll = cursorRow - textHeight + 1
	}

	scrollEnd := e.scroll + textHeight

	for i := e.scroll; i < scrollEnd; i++ {
		termPrintf(e.Bounds.Left + 1, e.Bounds.Top + 1 + i - e.scroll,
			   virtualLines[i])
	}

	RenderBorder(e.Bounds)

	if e.focus {
		termbox.SetCursor(e.Bounds.Left + 1 + cursorCol,
				  e.Bounds.Top  + 1 + cursorRow - e.scroll)
	}

	if e.mode == InsertMode {
		termPrintf(e.Bounds.Left + 1, e.Bounds.Bottom() - 1,
			   "-- INSERT --")
	}
}

func (e *Editbox) SetFocus() {
	e.focus = true
}

func (e *Editbox) UnsetFocus() {
	e.focus = false
}

func (e *Editbox) HandleEvent(ev escapebox.Event) {
	if ev.Type != termbox.EventKey {
		return
	}

	if e.mode == CommandMode {
		switch ev.Ch {
		case 'h':
			e.cursorChar--
		case 'l':
			e.cursorChar++
		case 'k':
			e.cursorLine--
		case 'j':
			e.cursorLine++
		case '0':
			e.cursorChar = 0
		case 'i':
			e.mode = InsertMode
		}
	} else if e.mode == InsertMode {
		if ev.Key == termbox.KeyEsc {
			e.mode = CommandMode
			e.cursorChar--
		} else if renderableChar(ev.Key) {
			line := e.lines[e.cursorLine]
			e.lines[e.cursorLine] =
				line[0:e.cursorChar] +
				string(ev.Ch) +
				line[e.cursorChar:len(line)]
			e.cursorChar++
		} else {
			switch (ev.Key) {
			case termbox.KeyArrowLeft:
				e.cursorChar--
			case termbox.KeyArrowRight:
				e.cursorChar++
			case termbox.KeyArrowUp:
				e.cursorLine--
			case termbox.KeyArrowDown:
				e.cursorLine++
			}
		}
	}

	e.cursorLine = max(0, e.cursorLine)
	e.cursorLine = min(len(e.lines) - 1, e.cursorLine)

	e.cursorChar = max(0, e.cursorChar)
	if e.mode == InsertMode {
		e.cursorChar = min(len(e.lines[e.cursorLine]), e.cursorChar)
	} else {
		e.cursorChar = min(len(e.lines[e.cursorLine]) - 1,
				   e.cursorChar)
	}
}
