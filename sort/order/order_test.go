package order_test

import (
	"encoding/json"
	"github.com/maurofran/page/sort/order"
	"github.com/maurofran/page/sort/order/direction"
	"github.com/maurofran/page/sort/order/null"
	"github.com/stretchr/testify/assert"
	"testing"
)

func makeOrder(direction direction.Direction, property string) order.Order {
	return order.New(direction, property)[0]
}

func TestNew(t *testing.T) {
	tests := map[string]struct {
		direction  direction.Direction
		properties []string
	}{
		"no_properties":    {direction.Desc, nil},
		"empty_properties": {direction.Desc, []string{}},
		"properties":       {direction.Desc, []string{"id", "name", "age"}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := order.New(test.direction, test.properties...)
			for i, o := range actual {
				assert.Equal(t, test.properties[i], o.Property())
				assert.Equal(t, test.direction, o.Direction())
				assert.False(t, o.IsIgnoreCase())
				assert.Equal(t, null.DefaultHandling, o.NullHandling())
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := map[string]struct {
		str      string
		expected order.Order
		err      error
	}{
		"empty_string":           {"", order.Order{}, order.ErrParse},
		"only_property":          {"id", makeOrder(direction.Default, "id"), nil},
		"property_direction":     {"id,desc", makeOrder(direction.Desc, "id"), nil},
		"property_ignore_case":   {"id,asc,ignore_case", makeOrder(direction.Asc, "id").IgnoreCase(), nil},
		"property_null_handling": {"id,asc,,nulls_first", makeOrder(direction.Asc, "id").NullsFirst(), nil},
		"too_many_parts":         {"id,asc,ignore_case,nulls_last,foo", makeOrder(direction.Asc, "id").IgnoreCase().NullsLast(), nil},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := order.Parse(test.str)
			assert.ErrorIs(t, err, test.err)
			assert.Equal(t, test.expected.Property(), actual.Property())
			assert.Equal(t, test.expected.Direction(), actual.Direction())
			assert.Equal(t, test.expected.IsIgnoreCase(), actual.IsIgnoreCase())
			assert.Equal(t, test.expected.NullHandling(), actual.NullHandling())
		})
	}
}

func TestBy(t *testing.T) {
	fixture := order.By("name")
	assert.Equal(t, "name", fixture.Property())
	assert.Equal(t, direction.Default, fixture.Direction())
	assert.False(t, fixture.IsIgnoreCase())
	assert.Equal(t, null.DefaultHandling, fixture.NullHandling())
}

func TestAsc(t *testing.T) {
	fixture := order.Asc("name")
	assert.Equal(t, "name", fixture.Property())
	assert.Equal(t, direction.Asc, fixture.Direction())
	assert.False(t, fixture.IsIgnoreCase())
	assert.Equal(t, null.DefaultHandling, fixture.NullHandling())
}

func TestDesc(t *testing.T) {
	fixture := order.Desc("name")
	assert.Equal(t, "name", fixture.Property())
	assert.Equal(t, direction.Desc, fixture.Direction())
	assert.False(t, fixture.IsIgnoreCase())
	assert.Equal(t, null.DefaultHandling, fixture.NullHandling())
}

func TestOrder_Property(t *testing.T) {
	fixture := makeOrder(direction.Asc, "name")
	assert.Equal(t, "name", fixture.Property())
}

func TestOrder_Direction(t *testing.T) {
	fixture := makeOrder(direction.Desc, "name")
	assert.Equal(t, direction.Desc, fixture.Direction())
}

func TestOrder_IsAscending(t *testing.T) {
	fixture := makeOrder(direction.Asc, "name")
	assert.True(t, fixture.IsAscending())
}

func TestOrder_IsDescending(t *testing.T) {
	fixture := makeOrder(direction.Desc, "name")
	assert.True(t, fixture.IsDescending())
}

func TestOrder_IsIgnoreCase(t *testing.T) {
	fixture := makeOrder(direction.Asc, "name").IgnoreCase()
	assert.True(t, fixture.IsIgnoreCase())
}

func TestOrder_NullHandling(t *testing.T) {
	fixture := makeOrder(direction.Asc, "name").NullsFirst()
	assert.Equal(t, null.First, fixture.NullHandling())
}

func TestOrder_With(t *testing.T) {
	fixture := makeOrder(direction.Asc, "name")
	actual := fixture.With(direction.Desc)
	assert.Equal(t, fixture.Property(), actual.Property())
	assert.Equal(t, direction.Desc, actual.Direction())
	assert.Equal(t, fixture.IsIgnoreCase(), actual.IsIgnoreCase())
	assert.Equal(t, fixture.NullHandling(), actual.NullHandling())
}

func TestOrder_Reverse(t *testing.T) {
	fixture := makeOrder(direction.Asc, "name")
	actual := fixture.Reverse()
	assert.Equal(t, fixture.Property(), actual.Property())
	assert.Equal(t, direction.Desc, actual.Direction())
	assert.Equal(t, fixture.IsIgnoreCase(), actual.IsIgnoreCase())
	assert.Equal(t, fixture.NullHandling(), actual.NullHandling())
}

func TestOrder_WithProperty(t *testing.T) {
	fixture := makeOrder(direction.Asc, "name")
	actual := fixture.WithProperty("age")
	assert.Equal(t, "age", actual.Property())
	assert.Equal(t, fixture.Direction(), actual.Direction())
	assert.Equal(t, fixture.IsIgnoreCase(), actual.IsIgnoreCase())
	assert.Equal(t, fixture.NullHandling(), actual.NullHandling())
}

func TestOrder_NullsNative(t *testing.T) {
	fixture := makeOrder(direction.Asc, "name").WithNullHandling(null.First)
	actual := fixture.NullsNative()
	assert.Equal(t, fixture.Property(), actual.Property())
	assert.Equal(t, fixture.Direction(), actual.Direction())
	assert.Equal(t, fixture.IsIgnoreCase(), actual.IsIgnoreCase())
	assert.Equal(t, null.Native, actual.NullHandling())
}

func TestOrder_NullsFirst(t *testing.T) {
	fixture := makeOrder(direction.Asc, "name")
	actual := fixture.NullsFirst()
	assert.Equal(t, fixture.Property(), actual.Property())
	assert.Equal(t, fixture.Direction(), actual.Direction())
	assert.Equal(t, fixture.IsIgnoreCase(), actual.IsIgnoreCase())
	assert.Equal(t, null.First, actual.NullHandling())
}

func TestOrder_NullsLast(t *testing.T) {
	fixture := makeOrder(direction.Asc, "name")
	actual := fixture.NullsLast()
	assert.Equal(t, fixture.Property(), actual.Property())
	assert.Equal(t, fixture.Direction(), actual.Direction())
	assert.Equal(t, fixture.IsIgnoreCase(), actual.IsIgnoreCase())
	assert.Equal(t, null.Last, actual.NullHandling())
}

func TestOrder_MarshalText(t *testing.T) {
	fixture := makeOrder(direction.Asc, "name").IgnoreCase().NullsLast()
	actual, err := json.Marshal(fixture)
	assert.NoError(t, err)
	assert.Equal(t, `"name,ASC,ignore_case,NULLS_LAST"`, string(actual))
}

func TestOrder_UnmarshalText(t *testing.T) {
	var actual order.Order
	err := json.Unmarshal([]byte(`"name,ASC,ignore_case,NULLS_LAST"`), &actual)
	assert.NoError(t, err)
	assert.Equal(t, "name", actual.Property())
	assert.Equal(t, direction.Asc, actual.Direction())
	assert.True(t, actual.IsIgnoreCase())
	assert.Equal(t, null.Last, actual.NullHandling())
}

func TestOrder_String(t *testing.T) {
	fixture := makeOrder(direction.Asc, "name").IgnoreCase().NullsLast()
	assert.Equal(t, "name: ASC, NULLS_LAST, ignoring case", fixture.String())
}
