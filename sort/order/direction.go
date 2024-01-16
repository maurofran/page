package order

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidDirection = errors.New("invalid direction")

// Direction is the enumeration for sort direction
type Direction int

const (
	DirectionAsc Direction = iota
	DirectionDesc
)

const defaultDirection = DirectionAsc

var directions = [...]string{
	"ASC",
	"DESC",
}

// ParseDirection parse a direction from a string.
func ParseDirection(value string) (Direction, error) {
	value = strings.ToUpper(value)
	for i, v := range directions {
		if v == value {
			return Direction(i), nil
		}
	}
	return DirectionAsc, fmt.Errorf("%w: value %q", ErrInvalidDirection, value)
}

// IsAscending returns whether the direction is ascending.
func (d Direction) IsAscending() bool {
	return d == DirectionAsc
}

// IsDescending returns whether the direction is descending.
func (d Direction) IsDescending() bool {
	return d == DirectionDesc
}

// Reverse the Direction.
func (d Direction) Reverse() Direction {
	if d == DirectionAsc {
		return DirectionDesc
	}
	return DirectionAsc
}

// MarshalText implements the encoding.TextMarshaler interface.
func (d Direction) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (d *Direction) UnmarshalText(text []byte) error {
	var err error
	*d, err = ParseDirection(string(text))
	return err
}

// String implements the fmt.Stringer interface.
func (d Direction) String() string {
	return directions[d]
}
