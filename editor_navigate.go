package main

import (
	"github.com/nsf/termbox-go"
)

type KeyCombo struct {
	Key termbox.Key
	Mod termbox.Modifier
	Ch  rune
}

var navigationCommands []map[KeyCombo]func()

func init() {
	navigationCommands = make([]map[KeyCombo]func(), 0, 2)

	defaultCommands := map[KeyCombo]func(){
		{termbox.KeyArrowLeft, 0, 0}: NavigateLeft,
		{termbox.KeyArrowRight, 0, 0}: NavigateRight,
		{termbox.KeyArrowUp, 0, 0}: NavigateUp,
		{termbox.KeyArrowDown, 0, 0}: NavigateDown,
		{termbox.KeyEnter, 0, 0}: EditCell,
		{0, 0, 'h'}: NavigateLeft,
		{0, 0, 'l'}: NavigateRight,
		{0, 0, 'j'}: NavigateUp,
		{0, 0, 'k'}: NavigateDown,
		{0, 0, ':'}: EnterCommandMode,
		{termbox.KeyEsc, 0, 0}: exitApplication,
	}
	PushNavigationCommands(defaultCommands)
}

func PushNavigationCommands(commands map[KeyCombo]func()) {
	navigationCommands = append(navigationCommands, commands)
}

func PopNavigationCommands() {

}

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

func handleNavigateMode(key termbox.Key, mod termbox.Modifier, ch rune) {
	keys := KeyCombo{key, mod, ch}

	for i := len(navigationCommands)-1; i >= 0; i-- {
		cmd, exists := navigationCommands[i][keys]
		if exists {
			cmd()
			break
		}
	}
}
