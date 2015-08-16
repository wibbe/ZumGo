package main

import "github.com/nsf/termbox-go"

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)
	termbox.HideCursor()

	clearScreen()
	termbox.Flush()

	initEditor()

	running := true

	for running {
		event := termbox.PollEvent()

		switch event.Type {
		case termbox.EventKey:
			if event.Key == termbox.KeyEsc {
				running = false
			}
		}

		redrawInterface()
		termbox.Flush()
	}

}
