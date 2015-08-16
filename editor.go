package main

import (
	"github.com/nsf/termbox-go"
)

type EditorMode int

const (
	NavigateMode EditorMode = iota
	EditMode
	CommandMode
	SearchMode
)

var currentDocument *Document
var currentMode EditorMode

func InitEditor() {
	currentMode = NavigateMode
	currentDocument = NewDocument()
}

func HandleKeyEvent(key termbox.Key, mod termbox.Modifier, ch rune) {
	switch currentMode {
	case NavigateMode:
		handleNavigateMode(key, mod, ch)
	case EditMode:
	case CommandMode:
	case SearchMode:
	}
}

func CurrentDoc() *Document {
	return currentDocument
}

func IsNavigationMode() bool {
	return currentMode == NavigateMode
}

func IsEditMode() bool {
	return currentMode == EditMode
}
