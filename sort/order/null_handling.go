package order

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidNullHandling = errors.New("invalid null handling")

// NullHandling is the enumeration for null handling hints that can be used in Order expressions.
type NullHandling int

const (
	NullHandlingNative NullHandling = iota
	NullHandlingNullsFirst
	NullHandlingNullsLast
)

const defaultNullHandling = NullHandlingNative

var nullHandling = [...]string{
	"NATIVE",
	"NULLS_FIRST",
	"NULLS_LAST",
}

// ParseNullHandling parse a NullHandling from a string.
func ParseNullHandling(value string) (NullHandling, error) {
	value = strings.ToUpper(value)
	for i, v := range nullHandling {
		if v == value {
			return NullHandling(i), nil
		}
	}
	return NullHandlingNative, fmt.Errorf("%w: value %q", ErrInvalidNullHandling, value)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (n NullHandling) MarshalText() ([]byte, error) {
	return []byte(n.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (n *NullHandling) UnmarshalText(text []byte) error {
	var err error
	*n, err = ParseNullHandling(string(text))
	return err
}

// String implements the fmt.Stringer interface.
func (n NullHandling) String() string {
	return nullHandling[n]
}
