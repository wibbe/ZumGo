package main

import (
	"errors"
	"strconv"
)

type Cell interface {
	String() string
	GetText() string
	Eval(doc *Document) (float64, error)
	Modified()
}

func NewCell(text string) Cell {
	if len(text) > 1 && text[0] == '=' {
		return &ExprCell{
			result:   0.0,
			modified: true,
			hasError: false,
		}
	} else {
		number, err := strconv.ParseFloat(text, 64)
		if err == nil {
			return &NumberCell{
				value: number,
			}
		}
	}
	return &TextCell{
		value: text,
	}
}

/// TextCell stores a plain-text string
type TextCell struct {
	value string
}

func (c *TextCell) String() string {
	return c.value
}

func (c *TextCell) GetText() string {
	return c.value
}

func (c *TextCell) Eval(doc *Document) (float64, error) {
	return 0.0, errors.New("Not allowed to evaluate text cell")
}

func (c *TextCell) Modified() {
}

/// NumberCell is used to store a number that can be used in cell evaluation
type NumberCell struct {
	value float64
}

func (c *NumberCell) String() string {
	return strconv.FormatFloat(c.value, 'f', -1, 64)
}

func (c *NumberCell) GetText() string {
	return c.String()
}

func (c *NumberCell) Eval(doc *Document) (float64, error) {
	return c.value, nil
}

func (c *NumberCell) Modified() {
}

/// ExprCell stores a mathematical experssion that can be evaluated
type ExprCell struct {
	result   float64
	modified bool
	hasError bool
}

func (c *ExprCell) String() string {
	if c.hasError {
		return "ERROR"
	}
	return strconv.FormatFloat(c.result, 'f', -1, 64)
}

func (c *ExprCell) GetText() string {
	return "="
}

func (c *ExprCell) Eval(doc *Document) (float64, error) {
	c.modified = false
	c.hasError = false
	c.result = 0.0

	return c.result, nil
}

func (c *ExprCell) Modified() {
	c.modified = true
}
