package page

import "github.com/maurofran/page/sort"

// Slice of data that indicates whether there's a next or previous slice available. Allows to obtain a
// Request to request a previous or next Slice.
type Slice[T any] interface {
	// Number returns the number of the current Slice. Is always non-negative.
	Number() uint
	// Size returns the size of the Slice.
	Size() uint
	// NumberOfElements returns the number of elements currently on this Slice.
	NumberOfElements() uint
	// Content returns the content of the Slice as a native slice.
	Content() []T
	// HasContent returns whether the Slice has content at all.
	HasContent() bool
	// Sort returns the sorting parameters for the Slice.
	Sort() *sort.Sort
	// IsFirst returns whether the current Slice is the first one.
	IsFirst() bool
	// IsLast returns whether the current Slice is the last one.
	IsLast() bool
	// HasNext returns if there is a next Slice.
	HasNext() bool
	// HasPrevious returns if there is a previous Slice.
	HasPrevious() bool
	// Request returns the Request that's been used to reques the current Slice.
	Request() Request // TODO PageRequest.of(Number(), Size(), Sort())
	// NextPageable returns the Request to request the next Slice. Can be Unpaged in case the current Slice is
	// already the last one. Clients should check HasNext before calling this method.
	NextPageable() Request
	// PreviousPageable returns the Request to request the previous Slice. Can be Unpaged in case the current Slice is
	// already the first one. Clients should check HasPrevious before calling this method.
	PreviousPageable() Request
	// NextOrLastPageable returns the Request describing the next Slice or the one describing the current slice in case
	// it's the last one.
	NextOrLastPageable() Request // TODO HasNext() ? NextPageable() : Request()
	// PreviousOrFirstPageable returns the Request describing the previous Slice or the one describing the current
	// slice in case it's the first one.
	PreviousOrFirstPageable() Request // TODO HasPrevious() ? PreviousPageable() : Request()
}
