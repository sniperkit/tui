package main

import (
	"github.com/nsf/termbox-go"
	"github.com/briansteffens/escapebox"
)

// Non-standard escape sequences
const (
	SeqShiftTab = 1
)

func buttonClickHandler(b *button) {
	panic("clicked!")
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc) // | termbox.InputMouse)

	escapebox.Init()
	defer escapebox.Close()

	escapebox.Register(SeqShiftTab, 91, 90)

	edit1 := Editbox {
		Bounds: Rect { Left: 2, Top: 6, Width: 30, Height: 10 },
		Value: "Hello!",
		cursor: 0,
		scroll: 0,
	}

	l := Label {
		Bounds: Rect { Left: 2, Top: 1, Width: 20, Height: 1 },
		Text: "Greetings:",
	}

	t := Textbox {
		Bounds: Rect { Left: 2, Top: 2, Width: 5, Height: 3 },
		Value: "12",
		cursor: 2,
		scroll: 0,
	}

	t2 := Textbox {
		Bounds: Rect { Left: 10, Top: 2, Width: 15, Height: 3},
		Value: "Greetings!",
		cursor: 0,
		scroll: 0,
	}

	checkbox1 := Checkbox {
		Bounds: Rect { Left: 27, Top: 1, Width: 30, Height: 1},
		Text: "Enable the whateverthing",
	}

	button1 := button {
		Bounds: Rect { Left: 27, Top: 2, Width: 10, Height: 3},
		Text: "Continue!",
		ClickHandler: buttonClickHandler,
	}

	c := Container {
		Controls: []Control {&t, &edit1, &l, &t2, &checkbox1,
				     &button1},
	}

	c.FocusNext()
	refresh(c)

	loop: for {
		ev := escapebox.PollEvent()

		handled := false

		switch ev.Seq {
		case escapebox.SeqNone:
			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyCtrlA:
					l.Text = ""
				case termbox.KeyCtrlC:
					break loop
				case termbox.KeyTab:
					c.FocusNext()
					handled = true
				}
			}
		case SeqShiftTab:
			c.FocusPrevious()
			handled = true
		}

		if !handled && c.Focused != nil {
			c.Focused.HandleEvent(ev)
		}

		refresh(c)
	}
}
