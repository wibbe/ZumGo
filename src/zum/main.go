package main

import (
	"app"
	"log"
	"os"
)

var zum *Zum = nil

func main() {
	err := app.Init("Zum", 800, 500)
	if err != nil {
		panic(err)
	}
	defer app.Shutdown()

	// Setup log
	logFile, err := os.Create(".zum.log")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Ltime | log.Lshortfile)
	log.Println("Application initialized")

	zum = NewZum()
	app.SetApplication(zum)

	InitEditor()

	app.Run()
}
