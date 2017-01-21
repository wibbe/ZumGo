package main

import (
	"log"
	"os"

	"github.com/nsf/termbox-go"
)

var applicationRunning bool

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)
	termbox.HideCursor()

	// Setup log
	logFile, err := os.Create(".zum.log")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Ltime | log.Lshortfile)

	log.Println("Application initialized")

	InitEditor()
	RedrawInterface()

	applicationRunning = true

	for applicationRunning {
		event := termbox.PollEvent()

		switch event.Type {
		case termbox.EventKey:
			HandleKeyEvent(event.Key, event.Mod, event.Ch)
		}

		RedrawInterface()
	}

}
