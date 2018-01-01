package main

import (
	"app"
	"log"
)

var applicationRunning bool

const (
	FontWidth  int = 9
	FontHeight     = 13
)

type Zum struct {
	width  int
	height int
	font   app.Font
}

func NewZum() *Zum {
	font := app.CreateFont("Consolas", 11.0, app.FontWeightNormal)
	if font == app.Font(-1) {
		panic("Could not create font")
	}

	return &Zum{
		width:  0,
		height: 0,
		font:   font,
	}
}

func (z *Zum) OnResize(w, h int) {
	log.Println("Resize")
	z.width = w
	z.height = h
}

func (z *Zum) OnPaint(w, h int) {
	z.width = w
	z.height = h
	app.FillRect(app.NewRectI(0, 0, w, h), app.WhiteBrush)
	RedrawInterface()
}

func (z *Zum) OnMouseMove(x, y float32) {

}

func (z *Zum) OnMouseEvent(button app.MouseButton, event app.MouseEvent, x, y float32) {

}

func (z *Zum) Size() (w, h int) {
	w = z.width/FontWidth + 1
	h = z.height/FontHeight + 1
	return
}

func (z *Zum) SetCell(col, row int, ch rune) {
	app.DrawChar(ch, z.font, app.BlackBrush, app.NewRectI(col*FontWidth, row*FontHeight, (col+1)*FontWidth, (row+1)*FontHeight), app.AlignVCenter|app.AlignHCenter)
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
