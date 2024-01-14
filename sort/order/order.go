package order

import (
	"fmt"
	"strings"
)

// Order implements the pairing of a Direction and a property. It is used to provide input for Sort.
type Order struct {
	property     string
	direction    Direction
	ignoreCase   bool
	nullHandling NullHandling
}

// New creates a new Order slice for provided Direction and properties.
func New(direction Direction, properties ...string) []Order {
	orders := make([]Order, len(properties))
	for i, property := range properties {
		orders[i] = Order{
			property:     property,
			direction:    direction,
			ignoreCase:   false,
			nullHandling: defaultNullHandling,
		}
	}
	return orders
}

// By creates a new Order instance with default sort direction.
func By(property string) Order {
	return Order{
		property:     property,
		direction:    defaultDirection,
		ignoreCase:   false,
		nullHandling: defaultNullHandling,
	}
}

// Asc creates a new Order instance with ascending sort direction.
func Asc(property string) Order {
	return Order{
		property:     property,
		direction:    DirectionAsc,
		ignoreCase:   false,
		nullHandling: defaultNullHandling,
	}
}

// Desc creates a new Order instance with descending sort direction.
func Desc(property string) Order {
	return Order{
		property:     property,
		direction:    DirectionDesc,
		ignoreCase:   false,
		nullHandling: defaultNullHandling,
	}
}

// Property returns the property to order for.
func (o Order) Property() string {
	return o.property
}

// Direction returns the direction to order by.
func (o Order) Direction() Direction {
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
func (o Order) NullHandling() NullHandling {
	return o.nullHandling
}

// With returns a new Order with given Direction
func (o Order) With(direction Direction) Order {
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
func (o Order) WithNullHandling(nullHandling NullHandling) Order {
	return Order{
		property:     o.property,
		direction:    o.direction,
		ignoreCase:   o.ignoreCase,
		nullHandling: nullHandling,
	}
}

// NullsFirst returns a new Order with NullHandlingNullsFirst NullHandling.
func (o Order) NullsFirst() Order {
	return o.WithNullHandling(NullHandlingNullsFirst)
}

// NullsLast returns a new Order with NullHandlingNullsLast NullHandling.
func (o Order) NullsLast() Order {
	return o.WithNullHandling(NullHandlingNullsLast)
}

// NullsNative returns a new Order with NullHandlingNative NullHandling.
func (o Order) NullsNative() Order {
	return o.WithNullHandling(NullHandlingNative)
}

func (o Order) String() string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("%s: %s", o.property, o.direction))
	if o.nullHandling != NullHandlingNative {
		result.WriteString(", ")
		result.WriteString(o.nullHandling.String())
	}
	if o.ignoreCase {
		result.WriteString(", ignoring case")
	}
	return result.String()
}
