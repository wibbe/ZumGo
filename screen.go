package main

import (
	"github.com/nsf/termbox-go"
	"unicode/utf8"
)

const (
	RowHeaderWidth = 8
)

type columnInfo struct {
	column int
	x      int
	width  int
}

func clearScreen() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func calculateColumnInfo(doc *Document) []columnInfo {
	info := make([]columnInfo, 0)

	w, _ := termbox.Size()

	x := RowHeaderWidth
	for i := doc.Scroll.X; i < doc.Width; i++ {
		width := doc.ColumnWidth[i]
		info = append(info, columnInfo{column: i, x: x, width: width})

		x += width
		if x > w {
			break
		}
	}

	return info
}

func drawText(x, y, length int, fg, bg termbox.Attribute, str string, format Format) {
	if length == -1 {
		for i, char := range str {
			termbox.SetCell(x+i, y, char, fg, bg)
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

			style := fg
			termbox.SetCell(x+i, y, ch, style, bg)
		}
	}
}

func getHeaderColor(cursor, header int) termbox.Attribute {
	color := termbox.ColorDefault

	if cursor == header {
		color = termbox.AttrReverse | termbox.ColorDefault
	}

	return color
}

func drawHeaders(doc *Document, info []columnInfo) {
	// Draw column headers
	for x := 0; x < len(info); x++ {
		color := getHeaderColor(doc.Cursor.X, info[x].column)
		drawText(info[x].x, 0, info[x].width, color, color, columnToStr(info[x].column), AlignCenter)
	}

	// Draw row headers
	_, h := termbox.Size()

	yEnd := h - 2
	if doc.Height < (h - 2) {
		yEnd = doc.Height
	}

	for y := 1; y <= yEnd; y++ {
		row := y + doc.Scroll.Y - 1
		color := getHeaderColor(doc.Cursor.Y, row)

		if row < doc.Height {
			drawText(0, y, RowHeaderWidth, color, color, rowToStr(row)+" ", AlignRight)
		}
	}
}

func drawDocument(doc *Document) {
	columnInfo := calculateColumnInfo(doc)
	drawHeaders(doc, columnInfo)
}

func redrawInterface() {
	clearScreen()

	if currentDocument != nil && currentDocument.Width > 0 && currentDocument.Height > 0 {
		drawDocument(currentDocument)
	}
}
