package main

import (
	"github.com/nsf/termbox-go"
)

type KeyCombo struct {
	Key termbox.Key
	Mod termbox.Modifier
	Ch  rune
}

var navigationCommands []map[KeyCombo]string

func init() {
	navigationCommands = make([]map[KeyCombo]string, 0, 2)

	defaultCommands := map[KeyCombo]string{
		{termbox.KeyArrowLeft, 0, 0}:  "navigate-left",
		{termbox.KeyArrowRight, 0, 0}: "navigate-right",
		{termbox.KeyArrowUp, 0, 0}:    "navigate-up",
		{termbox.KeyArrowDown, 0, 0}:  "navigate-down",
		{termbox.KeyEnter, 0, 0}:      "edit-current-cell",
		{0, 0, 'h'}:                   "navigate-left",
		{0, 0, 'l'}:                   "navigate-right",
		{0, 0, 'j'}:                   "navigate-up",
		{0, 0, 'k'}:                   "navigate-down",
		{0, 0, ':'}:                   "enter-command-mode",
		{termbox.KeyCtrlS, 0, 0}:      "save-document",
		{termbox.KeyEsc, 0, 0}:        "quit",
	}
	PushNavigationCommands(defaultCommands)
}

func PushNavigationCommands(commands map[KeyCombo]string) {
	navigationCommands = append(navigationCommands, commands)
}

func PopNavigationCommands() {

}

func handleNavigateMode(key termbox.Key, mod termbox.Modifier, ch rune) {
	keys := KeyCombo{key, mod, ch}

	for i := len(navigationCommands) - 1; i >= 0; i-- {
		cmd, exists := navigationCommands[i][keys]
		if exists {
			ExecLine(cmd)
			break
		}
	}
}
