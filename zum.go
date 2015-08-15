package main

import "github.com/nsf/termbox-go"

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)
	termbox.SetCursor(0, 0)

	clearScreen()
	termbox.Flush()

	running := true

	for running {
		event := termbox.PollEvent()

		switch event.Type {
		case termbox.EventKey:
			if event.Key == termbox.KeyEsc {
				running = false
			}
		}

		redrawScreen()
		termbox.Flush()
	}

}
