package page

// InitialOffsetScrollPosition is the initial scroll position to start scrolling using offset/limit.
const InitialOffsetScrollPosition = OffsetScrollPosition(0)

// OffsetScrollPosition is a ScrollPosition based on the offsets within a query.
type OffsetScrollPosition uint

// OffsetScrollPositionOf creates a new OffsetScrollPosition from an offset.
func OffsetScrollPositionOf(offset uint) OffsetScrollPosition {
	return OffsetScrollPositionOf(offset)
}

// Offset gets the zero or positive offset.
func (o OffsetScrollPosition) Offset() uint {
	return uint(o)
}

// AdvanceBy returns a new OffsetScrollPosition that has been advanced by the given value.
// Negative deltas will be constrained so that the new offset is at least zero.
func (o OffsetScrollPosition) AdvanceBy(delta int) OffsetScrollPosition {
	return OffsetScrollPosition(max(int(o)+delta, 0))
}

// IsInitial check if the OffsetScrollPosition is the initial value.
func (o OffsetScrollPosition) IsInitial() bool {
	return o == InitialOffsetScrollPosition
}
