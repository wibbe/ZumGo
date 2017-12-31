package main

import (
	"unicode/utf8"

	"github.com/nsf/termbox-go"
)

const (
	RowHeaderWidth = 8
)

var (
	ViewportSize Index
)

type columnInfo struct {
	column int
	x      int
	width  int
}

func clearScreen() {
	//termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func calculateColumnInfo(doc *Document) ([]columnInfo, int) {
	info := make([]columnInfo, 0)

	w, _ := zum.Size()

	x := RowHeaderWidth
	col := doc.Scroll.X

	for x < w {
		width := doc.GetColumnWidth(col)

		info = append(info, columnInfo{column: col, x: x, width: width})
		x += width
		col++
	}

	return info, col
}

func drawText(x, y, length int, fg, bg termbox.Attribute, str string, format Format) {
	if length == -1 {
		for i, char := range str {
			//termbox.SetCell(x+i, y, char, fg, bg)
			zum.SetCell(x+i, y, char)
		}
	} else {
		start := 0

		strLen := utf8.RuneCountInString(str)
		runes := make([]rune, strLen)
		for i, ch := range str {
			runes[i] = ch
		}

		switch format & AlignMask {
		case AlignLeft:
			start = 0
		case AlignCenter:
			start = (length / 2) - (strLen / 2)
		case AlignRight:
			start = length - strLen
		}

		for i := 0; i < length; i++ {
			charIdx := i - start
			ch := ' '
			if charIdx >= 0 && charIdx < strLen {
				ch = runes[charIdx]
			}

			//style := fg
			//termbox.SetCell(x+i, y, ch, style, bg)
			zum.SetCell(x+i, y, ch)
		}
	}
}

func getHeaderColor(equal bool) termbox.Attribute {
	color := termbox.ColorDefault

	if equal {
		color = termbox.AttrReverse | termbox.ColorDefault
	}

	return color
}

func drawHeaders(doc *Document, info []columnInfo) {
	// Draw column headers
	for x := 0; x < len(info); x++ {
		color := getHeaderColor(doc.Cursor.X == info[x].column)
		drawText(info[x].x, 0, info[x].width, color, color, columnToStr(info[x].column), AlignCenter)
	}

	// Draw row headers
	_, h := zum.Size()

	yEnd := h - 2
	if doc.Height < (h - 2) {
		yEnd = doc.Height
	}

	for y := 1; y <= yEnd; y++ {
		row := y + doc.Scroll.Y - 1
		color := getHeaderColor(doc.Cursor.Y == row)

		if row < doc.Height {
			drawText(0, y, RowHeaderWidth, color, color, rowToStr(row)+" ", AlignRight)
		}
	}
}

func drawWorkspace(doc *Document, info []columnInfo) {
	for row := doc.Scroll.Y; row < doc.Height; row++ {
		for i := 0; i < len(info); i++ {
			y := row - doc.Scroll.Y + 1

			color := getHeaderColor(doc.Cursor.X == info[i].column && doc.Cursor.Y == row)
			drawText(info[i].x, y, info[i].width, color, color, doc.GetCellDisplayText(NewIndex(info[i].column, row)), AlignLeft)
		}
	}
}

func drawDocument(doc *Document, columnInfo []columnInfo) {
	drawHeaders(doc, columnInfo)
	drawWorkspace(doc, columnInfo)
}

func drawFooter(doc *Document) {
	w, h := zum.Size()

	filename := "[No Name]"
	if doc.Filename != "" {
		filename = doc.Filename
	}

	if doc.Changed {
		filename = filename + "*"
	}

	cursorPos := " " + columnToStr(doc.Cursor.X) + rowToStr(doc.Cursor.Y) + " "

	color := termbox.ColorDefault | termbox.AttrReverse

	footerPos := h - 2

	inputPrompt := GetInputPrompt()
	inputLine := GetInputLine()
	promptLen := utf8.RuneCountInString(inputPrompt)

	if IsInputMode() {
		termbox.SetCursor(promptLen+GetInputCursor(), footerPos+1)
	} else {
		termbox.HideCursor()
		inputPrompt = ""
		inputLine = doc.GetCellDisplayText(doc.Cursor)
	}

	drawText(0, footerPos, w, color, color, filename, AlignLeft)
	drawText(w-8, footerPos, 8, color, color, cursorPos, AlignRight)

	drawText(0, footerPos+1, w, termbox.ColorDefault, termbox.ColorDefault, inputPrompt+inputLine, AlignLeft)
}

func RedrawInterface() {
	clearScreen()

	doc := CurrentDoc()
	columnInfo, width := calculateColumnInfo(doc)
	ViewportSize.X = width - 1

	drawDocument(doc, columnInfo)
	drawFooter(doc)
	//termbox.Flush()
}
