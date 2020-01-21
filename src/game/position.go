package game

import "fmt"

type position struct {
	row    int
	column int
}

func newPosition(row int, column int) position {
	return position{
		row:    row,
		column: column,
	}
}

func toString(position position) string {
	return fmt.Sprintf("[row: %d, column: %d]", position.row, position.column)
}
