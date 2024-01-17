package null_test

import (
	"encoding/json"
	"github.com/maurofran/page/sort/order/null"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseHandling(t *testing.T) {
	tests := map[string]struct {
		value    string
		expected null.Handling
		err      error
	}{
		"native":      {"native", null.Native, nil},
		"nulls_first": {"nulls_first", null.First, nil},
		"nulls_last":  {"nulls_last", null.Last, nil},
		"invalid":     {"invalid", null.Native, null.ErrInvalidHandling},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := null.ParseHandling(test.value)
			assert.ErrorIs(t, err, test.err)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestNullHandling_MarshalText(t *testing.T) {
	tests := map[string]struct {
		fixture  null.Handling
		expected string
	}{
		"native":      {null.Native, `"NATIVE"`},
		"nulls_first": {null.First, `"NULLS_FIRST"`},
		"nulls_last":  {null.Last, `"NULLS_LAST"`},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := json.Marshal(test.fixture)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, string(actual))
		})
	}
}

func TestNullHandling_UnmarshalText(t *testing.T) {
	tests := map[string]struct {
		fixture  string
		expected null.Handling
		err      error
	}{
		"native":      {`"NATIVE"`, null.Native, nil},
		"nulls_first": {`"NULLS_FIRST"`, null.First, nil},
		"nulls_last":  {`"NULLS_LAST"`, null.Last, nil},
		"invalid":     {`"INVALID"`, null.Native, null.ErrInvalidHandling},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var actual null.Handling
			err := json.Unmarshal([]byte(test.fixture), &actual)
			assert.ErrorIs(t, err, test.err)
			assert.Equal(t, test.expected, actual)
		})
	}
}
