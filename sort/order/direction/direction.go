package direction

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalid = errors.New("invalid direction")

// Direction is the enumeration for sort direction
type Direction int

const (
	Asc Direction = iota
	Desc
)

// Default represents the default sort direction.
const Default = Asc

var directions = [...]string{
	"ASC",
	"DESC",
}

// Parse a Direction from a string.
func Parse(value string) (Direction, error) {
	value = strings.ToUpper(value)
	for i, v := range directions {
		if v == value {
			return Direction(i), nil
		}
	}
	return Asc, fmt.Errorf("%w: value %q", ErrInvalid, value)
}

// IsAscending returns whether the direction is ascending.
func (d Direction) IsAscending() bool {
	return d == Asc
}

// IsDescending returns whether the direction is descending.
func (d Direction) IsDescending() bool {
	return d == Desc
}

// Reverse the Direction.
func (d Direction) Reverse() Direction {
	if d == Asc {
		return Desc
	}
	return Asc
}

// MarshalText implements the encoding.TextMarshaler interface.
func (d Direction) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (d *Direction) UnmarshalText(text []byte) error {
	var err error
	*d, err = Parse(string(text))
	return err
}

// String implements the fmt.Stringer interface.
func (d Direction) String() string {
	return directions[d]
}
