package tui

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/briansteffens/escapebox"
)

type CheckBox struct {
	Bounds     Rect
	Text       string
	Checked    bool
	focus      bool
}

func (c *CheckBox) Render() {
	checkContent := " "

	if c.Checked {
		checkContent = "X"
	}

	s := fmt.Sprintf("[%s] %s", checkContent, c.Text)

	count := min(len(s), c.Bounds.Width)
	termPrintf(c.Bounds.Left, c.Bounds.Top, s[0:count])

	if c.focus {
		termbox.SetCursor(c.Bounds.Left + 1, c.Bounds.Top)
	}
}

func (c *CheckBox) SetFocus() {
	c.focus = true
}

func (c *CheckBox) UnsetFocus() {
	c.focus = false
}

func (c *CheckBox) HandleEvent(ev escapebox.Event) bool {
	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeySpace:
			c.Checked = !c.Checked
			return true
		}
	}

	return false
}
