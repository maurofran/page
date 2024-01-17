package null

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidHandling = errors.New("invalid null handling")

// Handling is the enumeration for null handling hints that can be used in Order expressions.
type Handling int

const (
	HandlingNative Handling = iota
	HandlingNullsFirst
	HandlingNullsLast
)

const DefaultHandling = HandlingNative

var handlings = [...]string{
	"NATIVE",
	"NULLS_FIRST",
	"NULLS_LAST",
}

// ParseHandling parse a Handling from a string.
func ParseHandling(value string) (Handling, error) {
	value = strings.ToUpper(value)
	for i, v := range handlings {
		if v == value {
			return Handling(i), nil
		}
	}
	return HandlingNative, fmt.Errorf("%w: value %q", ErrInvalidHandling, value)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (n Handling) MarshalText() ([]byte, error) {
	return []byte(n.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (n *Handling) UnmarshalText(text []byte) error {
	var err error
	*n, err = ParseHandling(string(text))
	return err
}

// String implements the fmt.Stringer interface.
func (n Handling) String() string {
	return handlings[n]
}
