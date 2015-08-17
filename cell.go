package main

type Cell struct {
	value string
}

func NewCell(text string) *Cell {
	return &Cell{
		value: text,
	}
}

func (c *Cell) String() string {
	return c.value
}

func (c *Cell) SetText(text string) {
	c.value = text
}

func (c *Cell) GetText() string {
	return c.value
}
