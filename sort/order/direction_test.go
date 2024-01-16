package order_test

import (
	"encoding/json"
	"github.com/maurofran/page/sort/order"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseDirection(t *testing.T) {
	tests := map[string]struct {
		value    string
		expected order.Direction
		err      error
	}{
		"asc":     {"asc", order.DirectionAsc, nil},
		"desc":    {"desc", order.DirectionDesc, nil},
		"invalid": {"invalid", order.DirectionAsc, order.ErrInvalidDirection},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := order.ParseDirection(test.value)
			assert.ErrorIs(t, err, test.err)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDirection_IsAscending(t *testing.T) {
	tests := map[string]struct {
		fixture  order.Direction
		expected bool
	}{
		"asc":  {order.DirectionAsc, true},
		"desc": {order.DirectionDesc, false},
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
		fixture  order.Direction
		expected bool
	}{
		"asc":  {order.DirectionAsc, false},
		"desc": {order.DirectionDesc, true},
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
		fixture  order.Direction
		expected order.Direction
	}{
		"asc":  {order.DirectionAsc, order.DirectionDesc},
		"desc": {order.DirectionDesc, order.DirectionAsc},
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
		fixture  order.Direction
		expected string
	}{
		"asc":  {order.DirectionAsc, `"ASC"`},
		"desc": {order.DirectionDesc, `"DESC"`},
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
		expected order.Direction
		err      error
	}{
		"asc":     {`"ASC"`, order.DirectionAsc, nil},
		"desc":    {`"DESC"`, order.DirectionDesc, nil},
		"invalid": {`"INVALID"`, order.DirectionAsc, order.ErrInvalidDirection},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var actual order.Direction
			err := json.Unmarshal([]byte(test.fixture), &actual)
			assert.ErrorIs(t, err, test.err)
			assert.Equal(t, test.expected, actual)
		})
	}
}
