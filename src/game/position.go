package game

import "fmt"

type Position struct {
	Row    int
	Column int
}

func newPosition(row int, column int) Position {
	return Position{
		Row:    row,
		Column: column,
	}
}

func (pos Position) ToString() string {
	return fmt.Sprintf("[row: %d, column: %d]", pos.Row, pos.Column)
}
