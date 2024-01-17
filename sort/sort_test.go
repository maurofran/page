package sort_test

import (
	"encoding/json"
	"github.com/maurofran/page/sort"
	"github.com/maurofran/page/sort/order"
	"github.com/maurofran/page/sort/order/direction"
	"github.com/maurofran/page/sort/order/null"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func forceParse(t *testing.T, value ...string) *sort.Sort {
	t.Helper()
	s, err := sort.Parse(value...)
	require.NoError(t, err)
	return s
}

func TestUnsorted(t *testing.T) {
	fixture := sort.Unsorted()
	assert.Len(t, fixture.Orders(), 0)
}

func TestNew(t *testing.T) {
	test := map[string]struct {
		direction  direction.Direction
		properties []string
	}{
		"no_properties":       {direction.Asc, nil},
		"single_property":     {direction.Asc, []string{"name"}},
		"multiple_properties": {direction.Asc, []string{"name", "age"}},
	}

	for name, test := range test {
		t.Run(name, func(t *testing.T) {
			fixture := sort.New(test.direction, test.properties...)
			assert.Len(t, fixture.Orders(), len(test.properties))
			for i, o := range fixture.Orders() {
				assert.Equal(t, test.properties[i], o.Property())
				assert.Equal(t, test.direction, o.Direction())
				assert.False(t, o.IsIgnoreCase())
				assert.Equal(t, null.Native, o.NullHandling())
			}
		})
	}
}

func TestParse(t *testing.T) {
	test := map[string]struct {
		values []string
		orders []order.Order
		err    error
	}{
		"no_values":     {nil, nil, nil},
		"one_value":     {[]string{"name,asc"}, order.New(direction.Asc, "name"), nil},
		"two_values":    {[]string{"name,asc", "age,desc"}, order.New(direction.Asc, "name", "age"), nil},
		"invalid_value": {[]string{"name,asc", "age,desc", ""}, nil, order.ErrParse},
	}

	for name, test := range test {
		t.Run(name, func(t *testing.T) {
			fixture, err := sort.Parse(test.values...)
			assert.ErrorIs(t, err, test.err)
			if err == nil {
				assert.Len(t, fixture.Orders(), len(test.values))
			} else {
				assert.Nil(t, fixture)
			}
		})
	}
}

func TestBy(t *testing.T) {
	test := map[string]struct {
		direction  direction.Direction
		properties []string
	}{
		"no_properties":       {direction.Asc, nil},
		"single_property":     {direction.Asc, []string{"name"}},
		"multiple_properties": {direction.Asc, []string{"name", "age"}},
	}

	for name, test := range test {
		t.Run(name, func(t *testing.T) {
			fixture := sort.By(test.direction, test.properties...)
			assert.Len(t, fixture.Orders(), len(test.properties))
			for i, o := range fixture.Orders() {
				assert.Equal(t, test.properties[i], o.Property())
				assert.Equal(t, test.direction, o.Direction())
				assert.False(t, o.IsIgnoreCase())
				assert.Equal(t, null.Native, o.NullHandling())
			}
		})
	}
}

func TestSort_Ascending(t *testing.T) {
	fixture, _ := sort.Parse("name,asc", "age,desc")
	actual := fixture.Ascending()
	assert.NotSame(t, fixture, actual)
	assert.Len(t, actual.Orders(), len(fixture.Orders()))
	for _, o := range actual.Orders() {
		assert.True(t, o.IsAscending())
	}
}

func TestSort_Descending(t *testing.T) {
	fixture, _ := sort.Parse("name,asc", "age,desc")
	actual := fixture.Descending()
	assert.NotSame(t, fixture, actual)
	assert.Len(t, actual.Orders(), len(fixture.Orders()))
	for _, o := range actual.Orders() {
		assert.True(t, o.IsDescending())
	}
}

func TestSort_IsSorted(t *testing.T) {
	tests := map[string]struct {
		fixture  *sort.Sort
		expected bool
	}{
		"unsorted": {sort.Unsorted(), false},
		"sorted":   {sort.New(direction.Asc, "name"), true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.fixture.IsSorted())
		})
	}
}

func TestSort_IsUnsorted(t *testing.T) {
	tests := map[string]struct {
		fixture  *sort.Sort
		expected bool
	}{
		"unsorted": {sort.Unsorted(), true},
		"sorted":   {sort.New(direction.Asc, "name"), false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.fixture.IsUnsorted())
		})
	}
}

func TestSort_And(t *testing.T) {
	tests := map[string]struct {
		a        *sort.Sort
		b        *sort.Sort
		expected *sort.Sort
	}{
		"two_unsorted": {sort.Unsorted(), sort.Unsorted(), sort.Unsorted()},
		"two_sorted": {
			sort.New(direction.Asc, "name"),
			sort.New(direction.Desc, "age"),
			forceParse(t, "name,asc", "age,desc"),
		},
		"sorted_unsorted": {
			sort.New(direction.Asc, "name"),
			sort.Unsorted(),
			forceParse(t, "name,asc"),
		},
		"unsorted_sorted": {
			sort.Unsorted(),
			sort.New(direction.Asc, "name"),
			forceParse(t, "name,asc"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := test.a.And(test.b)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestSort_Reverse(t *testing.T) {
	fixture := forceParse(t, "name,asc", "age,desc")
	actual := fixture.Reverse()
	assert.NotSame(t, fixture, actual)
	assert.Equal(t, forceParse(t, "name,desc", "age,asc"), actual)
}

func TestSort_OrderFor(t *testing.T) {
	fixture := forceParse(t, "name,asc", "age,desc")
	o, ok := fixture.OrderFor("name")
	assert.True(t, ok)
	assert.Equal(t, "name", o.Property())
	assert.Equal(t, direction.Asc, o.Direction())
	o, ok = fixture.OrderFor("missing")
	assert.False(t, ok)
}

func TestSort_MarshalJSON(t *testing.T) {
	fixture := forceParse(t, "name,asc", "age,desc")
	data, err := json.Marshal(fixture)
	require.NoError(t, err)
	assert.Equal(t, `["name","age,DESC"]`, string(data))
}

func TestSort_UnmarshalJSON(t *testing.T) {
	fixture := `["name","age,DESC"]`
	var actual sort.Sort
	err := json.Unmarshal([]byte(fixture), &actual)
	assert.NoError(t, err)
	assert.Equal(t, forceParse(t, "name,asc", "age,desc"), &actual)
}

func TestSort_String(t *testing.T) {
	t.Run("unsorted", func(t *testing.T) {
		assert.Equal(t, "UNSORTED", sort.Unsorted().String())
	})

	t.Run("sorted", func(t *testing.T) {
		fixture := sort.New(direction.Asc, "name", "age")
		assert.Equal(t, "name: ASC, age: ASC", fixture.String())
	})
}
