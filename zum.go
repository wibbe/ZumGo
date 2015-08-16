package main

import "github.com/nsf/termbox-go"

var running bool

func exitApplication() {
	running = false
}

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

	InitEditor()

	running = true

	for running {
		event := termbox.PollEvent()

		switch event.Type {
		case termbox.EventKey:
			HandleKeyEvent(event.Key, event.Mod, event.Ch)
		}

		redrawInterface()
		termbox.Flush()
	}

}
