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
		{0, 0, 'h'}:                   "navigate-left",
		{termbox.KeyArrowRight, 0, 0}: "navigate-right",
		{0, 0, 'l'}:                   "navigate-right",
		{termbox.KeyArrowUp, 0, 0}:    "navigate-up",
		{0, 0, 'j'}:                   "navigate-up",
		{termbox.KeyArrowDown, 0, 0}:  "navigate-down",
		{0, 0, 'k'}:                   "navigate-down",
		{termbox.KeyEnter, 0, 0}:      "edit-current-cell",
		{0, 0, 'i'}:                   "edit-current-cell",
		{0, 0, ':'}:                   "enter-command-mode",
		{0, 0, '+'}:                   "modify-column-width 1",
		{0, 0, '-'}:                   "modify-column-width -1",
		{termbox.KeyCtrlS, 0, 0}:      "save-document",
		{termbox.KeyCtrlO, 0, 0}:      "open-document",
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
