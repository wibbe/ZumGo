package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

const (
	DefaultWidth       = 5
	DefaultHeight      = 10
	DefaultColumnWidth = 20
	MinColumnWidth     = 3
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
	Cells       map[Index]Cell
	ColumnWidth []int
	Cursor      Index
	Scroll      Index
	Filename    string
	Changed     bool
}

func NewDocument() *Document {
	doc := &Document{
		Width:       DefaultWidth,
		Height:      DefaultHeight,
		Cells:       make(map[Index]Cell),
		ColumnWidth: make([]int, DefaultWidth),
		Scroll:      NewIndex(0, 0),
		Cursor:      NewIndex(0, 0),
		Filename:    "",
		Changed:     false,
	}

	for i := 0; i < doc.Width; i++ {
		doc.ColumnWidth[i] = DefaultColumnWidth
	}

	return doc
}

func LoadDocument(filename string) (*Document, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	var rows [][]string

	rows, err = csvReader.ReadAll()

	height := len(rows)
	if height == 0 {
		return nil, errors.New("No data in the document")
	}

	doc := &Document{
		Width:       0,
		Height:      0,
		Cells:       make(map[Index]Cell),
		ColumnWidth: nil,
		Scroll:      NewIndex(0, 0),
		Cursor:      NewIndex(0, 0),
		Filename:    filename,
		Changed:     false,
	}

	headerInfo := true
	headerRe := regexp.MustCompile("#{([0-9]+)}")

	// Parse header information
	for i := 0; i < len(rows[0]); i++ {
		matches := headerRe.FindStringSubmatch(rows[0][i])

		if len(matches) != 2 {
			headerInfo = false
		}

		if headerInfo {
			width, err := strconv.Atoi(matches[1])
			if err != nil {
				headerInfo = false
			} else {
				doc.ColumnWidth = append(doc.ColumnWidth, width)
			}
		}

		if !headerInfo {
			doc.ColumnWidth = append(doc.ColumnWidth, DefaultColumnWidth)
		}
	}

	startRow := 0
	if headerInfo {
		startRow = 1
	}

	doc.Height = height - startRow
	doc.Width = len(rows[0])

	for y := startRow; y < height; y++ {
		for x := 0; x < doc.Width; x++ {
			if rows[y][x] != "" {
				doc.Cells[NewIndex(x, y-startRow)] = NewCell(rows[y][x])
			}
		}
	}

	doc.Evaluate()
	log.Printf("Document '%s' loaded", filename)

	return doc, nil
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
	d.Cells[idx] = NewCell(text)
	d.Changed = true
	d.Evaluate()
}

func (d *Document) GetCell(idx Index) (Cell, error) {
	cell, exists := d.Cells[idx]
	if exists {
		return cell, nil
	}
	return nil, errors.New("No data at specified index")
}

func (d *Document) Evaluate() {
	for _, cell := range d.Cells {
		cell.Modified()
	}

	for _, cell := range d.Cells {
		cell.Eval(d)
	}
}

func (d *Document) ModifyColumnWidth(column, modification int) {
	if column < 0 || column >= d.Width {
		return
	}

	d.Changed = true
	d.ColumnWidth[column] += modification

	if d.ColumnWidth[column] < MinColumnWidth {
		d.ColumnWidth[column] = MinColumnWidth
	}
}

func (d *Document) Save() error {
	if d.Filename == "" {
		return errors.New("Could not save document, no filename specified")
	}

	file, err := os.Create(d.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)

	row := make([]string, d.Width)

	// Write column info
	for i := 0; i < d.Width; i++ {
		row[i] = fmt.Sprintf("#{%d}", d.ColumnWidth[i])
	}

	err = csvWriter.Write(row)
	if err != nil {
		return err
	}

	for y := 0; y < d.Height; y++ {
		for x := 0; x < d.Width; x++ {
			row[x] = d.GetCellText(NewIndex(x, y))
		}

		err = csvWriter.Write(row)
		if err != nil {
			return err
		}
	}

	csvWriter.Flush()
	err = csvWriter.Error()
	if err != nil {
		return err
	}

	log.Printf("Document '%s' saved", d.Filename)

	d.Changed = false

	return nil
}
