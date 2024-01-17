package sort

import (
	"encoding/json"
	"github.com/maurofran/page/sort/order"
	"github.com/maurofran/page/sort/order/direction"
	"strings"
)

// Sort option for queries. You have to provide at least a list of properties to sort for that must not include
// empty strings. The direction defaults to Asc.
type Sort struct {
	orders []order.Order
}

// Unsorted gets an unsorted sort specification
func Unsorted() *Sort {
	return &Sort{}
}

// New creates a new Sort instance with supplied direction and properties.
func New(direction direction.Direction, properties ...string) *Sort {
	if len(properties) == 0 {
		return Unsorted()
	}
	return &Sort{orders: order.New(direction, properties...)}
}

// Parse zero or more sort clauses for a Sort instance.
func Parse(values ...string) (*Sort, error) {
	if len(values) == 0 {
		return Unsorted(), nil
	}
	orders := make([]order.Order, len(values))
	for i, str := range values {
		o, err := order.Parse(str)
		if err != nil {
			return nil, err
		}
		orders[i] = o
	}
	return ByOrder(orders...), nil
}

// By creates a new Sort with provided properties.
func By(direction direction.Direction, properties ...string) *Sort {
	return New(direction, properties...)
}

// ByOrder creates a new Sort with provided Order.
func ByOrder(orders ...order.Order) *Sort {
	return &Sort{orders: orders}
}

// Descending returns a new Sort with the current setup but descending order direction.
func (s *Sort) Descending() *Sort {
	return s.withDirection(direction.Desc)
}

// Ascending returns a new Sort with the current setup but ascending order direction.
func (s *Sort) Ascending() *Sort {
	return s.withDirection(direction.Asc)
}

func (s *Sort) withDirection(direction direction.Direction) *Sort {
	orders := make([]order.Order, len(s.orders))
	for i, o := range s.orders {
		orders[i] = o.With(direction)
	}
	return ByOrder(orders...)
}

// IsSorted returns true if the Sort is sorted.
func (s *Sort) IsSorted() bool {
	return !s.IsEmpty()
}

// IsUnsorted returns true if the Sort is not sorted.
func (s *Sort) IsUnsorted() bool {
	return !s.IsSorted()
}

// IsEmpty returns true if the Sort is empty.
func (s *Sort) IsEmpty() bool {
	return s == nil || len(s.orders) == 0
}

// And returns a new Sort consisting of the order.Order of the current Sort combined with the given one.
func (s *Sort) And(other *Sort) *Sort {
	var orders []order.Order
	orders = append(orders, s.orders...)
	orders = append(orders, other.orders...)
	return ByOrder(orders...)
}

// Reverse returns a new Sort with reversed sort order.Order turning ascending into descending and vice versa.
func (s *Sort) Reverse() *Sort {
	orders := make([]order.Order, len(s.orders))
	for i, o := range s.orders {
		orders[i] = o.Reverse()
	}
	return ByOrder(orders...)
}

// OrderFor gets the order.Order clause for provided property.
func (s *Sort) OrderFor(property string) (order.Order, bool) {
	for _, o := range s.orders {
		if o.Property() == property {
			return o, true
		}
	}
	return order.Order{}, false
}

// Orders gets the slice of orders.
func (s *Sort) Orders() []order.Order {
	return s.orders
}

// MarshalJSON implements the json.Marshaler interface.
func (s *Sort) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.orders)
}

// UnmarshalJSON implements the json.Unmarshaler interface,
func (s *Sort) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &s.orders)
}

// String implements the fmt.Stringer interface.
func (s *Sort) String() string {
	if s.IsUnsorted() {
		return "UNSORTED"
	}
	parts := make([]string, len(s.orders))
	for i, o := range s.orders {
		parts[i] = o.String()
	}
	return strings.Join(parts, ", ")
}
