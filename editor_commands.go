package main

import (
	"log"
	"strings"
)

func NavigateUp() {
	doc := CurrentDoc()
	if doc.Cursor.Y > 0 {
		doc.Cursor.Y -= 1
	}
}

func NavigateDown() {
	doc := CurrentDoc()
	if doc.Cursor.Y < (doc.Height - 1) {
		doc.Cursor.Y += 1
	}
}

func NavigateLeft() {
	doc := CurrentDoc()
	if doc.Cursor.X > 0 {
		doc.Cursor.X -= 1
	}
}

func NavigateRight() {
	doc := CurrentDoc()
	if doc.Cursor.X < (doc.Width - 1) {
		doc.Cursor.X += 1
	}
}

func NavigateRightOrNewLine() {
	doc := CurrentDoc()
	if doc.Cursor.X < (doc.Width - 1) {
		doc.Cursor.X += 1
	} else if doc.Cursor.Y < (doc.Height - 1) {
		doc.Cursor.X = 0
		doc.Cursor.Y += 1
	}
}

func exitApplication() {
	quit, exists := GetArg(1)
	quit = strings.ToLower(quit)

	if exists && (quit == "y" || quit == "yes") {
		applicationRunning = false
	} else {
		EnableInputMode("Quit (Y/n): ", "", func(line string) {
			result := strings.ToLower(line)
			if result == "" || result == "y" || result == "yes" {
				applicationRunning = false
			}
		})
	}
}

func saveDocument() {
	doc := CurrentDoc()
	filename, exists := GetArg(1)

	if exists {
		doc.Filename = filename
	}

	if doc.Filename == "" {
		EnableInputMode("Save file: ", "", func(line string) {
			doc.Filename = line
			if doc.Filename != "" {
				err := doc.Save()
				if err != nil {
					log.Println("Could not save document: ", err)
				}
			}
		})
	} else {
		err := doc.Save()
		if err != nil {
			log.Println("Could not save document: ", err)
		}
	}
}
