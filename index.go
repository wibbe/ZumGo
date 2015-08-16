package main

import (
	"bytes"
	"fmt"
	"math"
)

// Index acts as a reference to cells in a document
type Index struct {
	X int
	Y int
}

func NewIndex(x, y int) Index {
	return Index{X: x, Y: y}
}

func rowToStr(row int) string {
	return fmt.Sprintf("%d", row+1)
}

func columnToStr(col int) string {
	power := 0
	for i := 0; i < 10; i++ {
		value := int(math.Pow(26.0, float64(power)))
		if value > col {
			power = i
			break
		}
	}

	var buffer bytes.Buffer

	for power >= 0 {
		value := int(math.Pow(26.0, float64(power)))
		if value > col {
			power -= 1
			if power < 0 {
				buffer.WriteRune('A')
			}
		} else {
			division := col / value
			remainder := col % value

			buffer.WriteRune(rune(int('A') + division))

			col = remainder
			power -= 1
		}
	}

	return buffer.String()
}
