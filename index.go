package main

// Index acts as a reference to cells in a document
type Index struct {
	X int
	Y int
}

func NewIndex(x, y int) Index {
	return Index{X: x, Y: y}
}
