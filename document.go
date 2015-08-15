package main

const (
	DefaultWidth  = 5
	DefaultHeight = 5
)

type Cell struct {
	value string
}

type Document struct {
	width  int
	height int
	cells  map[Index]*Cell
	cursor Index
}

func NewDocument() *Document {
	return &Document{width: DefaultWidth,
		height: DefaultHeight,
		cells:  make(map[Index]*Cell),
		cursor: NewIndex(0, 0)}
}
