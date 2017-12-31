package main

import (
	//	"github.com/nsf/termbox-go"
	"app"
	"log"
)

var applicationRunning bool

type Zum struct {
}

func NewZum() *Zum {
	return &Zum{}
}

func (z *Zum) OnResize(w, h int) {
	log.Println("Resize")
}

func (z *Zum) OnPaint(w, h int) {
	app.FillRect(app.NewRectI(0, w, 0, h), app.WhiteBrush)
}

func (z *Zum) OnMouseMove(x, y float32) {

}

func (z *Zum) OnMouseEvent(button app.MouseButton, event app.MouseEvent, x, y float32) {

}

/*
func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)
	termbox.HideCursor()


	InitEditor()
	RedrawInterface()

	applicationRunning = true

	for applicationRunning {
		event := termbox.PollEvent()

		switch event.Type {
		case termbox.EventKey:
			HandleKeyEvent(event.Key, event.Mod, event.Ch)
		}

		RedrawInterface()
	}

}
*/
