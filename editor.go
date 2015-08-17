package main

import (
	"github.com/nsf/termbox-go"
)

type EditorMode int

const (
	NavigateMode EditorMode = iota
	InputMode
)

var currentDocument *Document
var currentMode EditorMode

func InitEditor() {
	currentMode = NavigateMode
	currentDocument = NewDocument()

	RegisterCommands("default", []Command{
		{"navigate-left", NavigateLeft},
		{"navigate-right", NavigateRight},
		{"navigate-up", NavigateUp},
		{"navigate-down", NavigateDown},
		{"enter-command-mode", EnterCommandMode},
		{"edit-current-cell", EditCell},
		{"quit", exitApplication},
	})
}

func HandleKeyEvent(key termbox.Key, mod termbox.Modifier, ch rune) {
	switch currentMode {
	case NavigateMode:
		handleNavigateMode(key, mod, ch)
	case InputMode:
		handleInputMode(key, mod, ch)
	}
}

func CurrentDoc() *Document {
	return currentDocument
}

func IsNavigationMode() bool {
	return currentMode == NavigateMode
}

func IsInputMode() bool {
	return currentMode == InputMode
}

func EditCell() {
	EnableInputMode("", currentDocument.GetCellText(currentDocument.Cursor), cellEditFinished)
}

func cellEditFinished(line string) {
	currentDocument.SetCellText(currentDocument.Cursor, line)
	NavigateRightOrNewLine()
}
