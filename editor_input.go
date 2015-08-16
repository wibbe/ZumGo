package main

import (
	"github.com/nsf/termbox-go"
)

type InputCmd func(line string)

var inputPrompt string
var inputFunction InputCmd
var inputLine []rune
var inputCursor int

func EnableInputMode(prompt, value string, cmd InputCmd) {
	currentMode = InputMode
	inputPrompt = prompt
	inputFunction = cmd
	inputLine = make([]rune, 0, 128)

	for _, ch := range value {
		inputLine = append(inputLine, ch)
	}

	inputCursor = len(inputLine)
}

func DisableInputMode() {
	currentMode = NavigateMode
	inputPrompt = ""
	inputFunction = nil
	inputCursor = 0
}

func insertCh(ch rune) {
	inputLine = append(inputLine, ' ')
	copy(inputLine[inputCursor+1:], inputLine[inputCursor:])
	inputLine[inputCursor] = ch
	inputCursor += 1
}

func handleInputMode(key termbox.Key, mod termbox.Modifier, ch rune) {
	if ch == 0 {
		switch key {
		case termbox.KeyArrowLeft:
			if inputCursor > 0 {
				inputCursor -= 1
			}

		case termbox.KeyArrowRight:
			if inputCursor < len(inputLine) {
				inputCursor += 1
			}

		case termbox.KeyEnter:
			inputFunction(string(inputLine))
			DisableInputMode()

		case termbox.KeyEsc:
			DisableInputMode()

		case termbox.KeySpace:
			insertCh(' ')

		case termbox.KeyBackspace:
			if inputCursor > 0 {
				inputLine = append(inputLine[:inputCursor-1], inputLine[inputCursor:]...)
				inputCursor -= 1
			}

		case termbox.KeyDelete:
			if inputCursor < len(inputLine) {
				inputLine = append(inputLine[:inputCursor], inputLine[inputCursor+1:]...)
			}

		case termbox.KeyHome:
			inputCursor = 0

		case termbox.KeyEnd:
			inputCursor = len(inputLine)
		}
	} else {
		insertCh(ch)
	}
}

func GetInputLine() string {
	return string(inputLine)
}

func GetInputPrompt() string {
	return inputPrompt
}

func GetInputCursor() int {
	return inputCursor
}
