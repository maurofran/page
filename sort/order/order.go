package order

import (
	"errors"
	"fmt"
	"github.com/maurofran/page/sort/order/direction"
	"github.com/maurofran/page/sort/order/null"
	"strings"
)

// ErrParse is the error returned when an order cannot be parsed.
var ErrParse = errors.New("order parse error")

// Order implements the pairing of a Direction and a property. It is used to provide input for Sort.
type Order struct {
	property     string
	direction    direction.Direction
	ignoreCase   bool
	nullHandling null.Handling
}

// New creates a new Order slice for provided Direction and properties.
func New(direction direction.Direction, properties ...string) []Order {
	orders := make([]Order, len(properties))
	for i, property := range properties {
		orders[i] = Order{
			property:     property,
			direction:    direction,
			ignoreCase:   false,
			nullHandling: null.DefaultHandling,
		}
	}
	return orders
}

// Parse an Order formatted string
func Parse(data string) (Order, error) {
	var err error
	var order Order
	data = strings.TrimSpace(data)
	if data == "" {
		return order, fmt.Errorf("%w: empty string", ErrParse)
	}
	parts := strings.Split(data, ",")
	count := len(parts)
	if count >= 4 {
		order.nullHandling, err = null.ParseHandling(parts[3])
	}
	if count >= 3 {
		order.ignoreCase = parts[2] == "ignore_case"
	}
	if count >= 2 {
		order.direction, err = direction.Parse(parts[1])
	}
	if count >= 1 {
		order.property = parts[0]
	}
	return order, err
}

// By creates a new Order instance with default sort direction.
func By(property string) Order {
	return Order{
		property:     property,
		direction:    direction.Default,
		ignoreCase:   false,
		nullHandling: null.DefaultHandling,
	}
}

// Asc creates a new Order instance with ascending sort direction.
func Asc(property string) Order {
	return Order{
		property:     property,
		direction:    direction.Asc,
		ignoreCase:   false,
		nullHandling: null.DefaultHandling,
	}
}

// Desc creates a new Order instance with descending sort direction.
func Desc(property string) Order {
	return Order{
		property:     property,
		direction:    direction.Desc,
		ignoreCase:   false,
		nullHandling: null.DefaultHandling,
	}
}

// Property returns the property to order for.
func (o Order) Property() string {
	return o.property
}

// Direction returns the direction to order by.
func (o Order) Direction() direction.Direction {
	return o.direction
}

// IsAscending returns whether the sorting for the property is ascending.
func (o Order) IsAscending() bool {
	return o.direction.IsAscending()
}

// IsDescending returns whether the sorting for the property is descending.
func (o Order) IsDescending() bool {
	return o.direction.IsDescending()
}

// IsIgnoreCase returns whether the sorting should be case-insensitive.
func (o Order) IsIgnoreCase() bool {
	return o.ignoreCase
}

// NullHandling returns the null handling.
func (o Order) NullHandling() null.Handling {
	return o.nullHandling
}

// With returns a new Order with given Direction
func (o Order) With(direction direction.Direction) Order {
	return Order{
		property:     o.property,
		direction:    direction,
		ignoreCase:   o.ignoreCase,
		nullHandling: o.nullHandling,
	}
}

// Reverse the direction.
func (o Order) Reverse() Order {
	return Order{
		property:     o.property,
		direction:    o.direction.Reverse(),
		ignoreCase:   o.ignoreCase,
		nullHandling: o.nullHandling,
	}
}

// WithProperty returns a new Order with supplied property.
func (o Order) WithProperty(property string) Order {
	return Order{
		property:     property,
		direction:    o.direction,
		ignoreCase:   o.ignoreCase,
		nullHandling: o.nullHandling,
	}
}

// IgnoreCase returns a new Order with ignore case flag enabled.
func (o Order) IgnoreCase() Order {
	return Order{
		property:     o.property,
		direction:    o.direction,
		ignoreCase:   true,
		nullHandling: o.nullHandling,
	}
}

// WithNullHandling returns a new Order with provided NullHandling.
func (o Order) WithNullHandling(nullHandling null.Handling) Order {
	return Order{
		property:     o.property,
		direction:    o.direction,
		ignoreCase:   o.ignoreCase,
		nullHandling: nullHandling,
	}
}

// NullsFirst returns a new Order with First NullHandling.
func (o Order) NullsFirst() Order {
	return o.WithNullHandling(null.First)
}

// NullsLast returns a new Order with Last NullHandling.
func (o Order) NullsLast() Order {
	return o.WithNullHandling(null.Last)
}

// NullsNative returns a new Order with Native NullHandling.
func (o Order) NullsNative() Order {
	return o.WithNullHandling(null.Native)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (o Order) MarshalText() ([]byte, error) {
	var result strings.Builder
	result.WriteString(o.property)
	if o.direction != direction.Asc || o.ignoreCase || o.nullHandling != null.Native {
		result.WriteString(",")
		result.WriteString(o.direction.String())
	}
	if o.ignoreCase || o.nullHandling != null.Native {
		result.WriteString(",ignore_case")
	}
	if o.nullHandling != null.Native {
		result.WriteString(",")
		result.WriteString(o.nullHandling.String())
	}
	return []byte(result.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (o *Order) UnmarshalText(data []byte) error {
	var err error
	*o, err = Parse(string(data))
	return err
}

func (o Order) String() string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("%s: %s", o.property, o.direction))
	if o.nullHandling != null.Native {
		result.WriteString(", ")
		result.WriteString(o.nullHandling.String())
	}
	if o.ignoreCase {
		result.WriteString(", ignoring case")
	}
	return result.String()
}
