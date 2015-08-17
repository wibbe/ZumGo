package main

import (
	"github.com/nsf/termbox-go"
	"strings"
)

var running bool

func exitApplication() {
	quit, exists := GetArg(1)
	quit = strings.ToLower(quit)

	if exists && (quit == "y" || quit == "yes") {
		running = false
	} else {
		EnableInputMode("Quit (Y/n): ", "", func(line string) {
			result := strings.ToLower(line)
			if result == "" || result == "y" || result == "yes" {
				running = false
			}
		})
	}
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
