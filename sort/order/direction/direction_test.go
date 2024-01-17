package direction_test

import (
	"encoding/json"
	"github.com/maurofran/page/sort/order/direction"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseDirection(t *testing.T) {
	tests := map[string]struct {
		value    string
		expected direction.Direction
		err      error
	}{
		"asc":     {"asc", direction.Asc, nil},
		"desc":    {"desc", direction.Desc, nil},
		"invalid": {"invalid", direction.Asc, direction.ErrInvalid},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := direction.Parse(test.value)
			assert.ErrorIs(t, err, test.err)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDirection_IsAscending(t *testing.T) {
	tests := map[string]struct {
		fixture  direction.Direction
		expected bool
	}{
		"asc":  {direction.Asc, true},
		"desc": {direction.Desc, false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := test.fixture.IsAscending()
			assert.Equal(t, test.expected, actual)
		})
	}
}
func TestDirection_IsDescending(t *testing.T) {
	tests := map[string]struct {
		fixture  direction.Direction
		expected bool
	}{
		"asc":  {direction.Asc, false},
		"desc": {direction.Desc, true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := test.fixture.IsDescending()
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDirection_Reverse(t *testing.T) {
	tests := map[string]struct {
		fixture  direction.Direction
		expected direction.Direction
	}{
		"asc":  {direction.Asc, direction.Desc},
		"desc": {direction.Desc, direction.Asc},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := test.fixture.Reverse()
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDirection_MarshalText(t *testing.T) {
	tests := map[string]struct {
		fixture  direction.Direction
		expected string
	}{
		"asc":  {direction.Asc, `"ASC"`},
		"desc": {direction.Desc, `"DESC"`},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := json.Marshal(test.fixture)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, string(actual))
		})
	}
}

func TestDirection_UnmarshalText(t *testing.T) {
	tests := map[string]struct {
		fixture  string
		expected direction.Direction
		err      error
	}{
		"asc":     {`"ASC"`, direction.Asc, nil},
		"desc":    {`"DESC"`, direction.Desc, nil},
		"invalid": {`"INVALID"`, direction.Asc, direction.ErrInvalid},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var actual direction.Direction
			err := json.Unmarshal([]byte(test.fixture), &actual)
			assert.ErrorIs(t, err, test.err)
			assert.Equal(t, test.expected, actual)
		})
	}
}
