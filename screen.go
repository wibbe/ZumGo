package main

import "github.com/nsf/termbox-go"

func clearScreen() {
	w, h := termbox.Size()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			termbox.SetCell(x, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
		}
	}
}

func drawDocument(doc *Document) {

}

func redrawScreen() {
	clearScreen()
}
