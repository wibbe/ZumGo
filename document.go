package main

const (
	DefaultWidth       = 5
	DefaultHeight      = 10
	DefaultColumnWidth = 20
)

type Format uint16

const (
	AlignMask   Format = 0x000F
	AlignLeft          = 0x0001
	AlignCenter        = 0x0002
	AlignRight         = 0x0004
)

type Document struct {
	Width       int
	Height      int
	Cells       map[Index]*Cell
	ColumnWidth []int
	Cursor      Index
	Scroll      Index
	Filename    string
}

func NewDocument() *Document {
	doc := &Document{
		Width:       DefaultWidth,
		Height:      DefaultHeight,
		Cells:       make(map[Index]*Cell),
		ColumnWidth: make([]int, DefaultWidth),
		Scroll:      NewIndex(0, 0),
		Cursor:      NewIndex(0, 0),
		Filename:    "",
	}

	for i := 0; i < doc.Width; i++ {
		doc.ColumnWidth[i] = DefaultColumnWidth
	}

	return doc
}

func (d *Document) GetCellDisplayText(idx Index) string {
	cell, exists := d.Cells[idx]
	if exists && cell != nil {
		return cell.String()
	}
	return ""
}

func (d *Document) GetCellText(idx Index) string {
	cell, exists := d.Cells[idx]
	if exists && cell != nil {
		return cell.GetText()
	}
	return ""
}

func (d *Document) SetCellText(idx Index, text string) {
	cell, exists := d.Cells[idx]
	if exists {
		cell.SetText(text)
	} else {
		d.Cells[idx] = &Cell{value: text}
	}
}
