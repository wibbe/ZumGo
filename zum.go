package main

import (
	"github.com/nsf/termbox-go"
	"log"
	"os"
)

var applicationRunning bool

type ApplicationError struct {
	message string
}

func (a *ApplicationError) Error() string {
	return a.message
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

	// Setup log
	logFile, err := os.Create("zum.log")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	InitEditor()

	log.Println("Application initialized")

	applicationRunning = true

	for applicationRunning {
		event := termbox.PollEvent()

		switch event.Type {
		case termbox.EventKey:
			HandleKeyEvent(event.Key, event.Mod, event.Ch)
		}

		redrawInterface()
		termbox.Flush()
	}

}
