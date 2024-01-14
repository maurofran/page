package page

import "errors"

// ErrUnlimited is returned when a Max operation is invoked on an Unlimited Limit.
var ErrUnlimited = errors.New("limit is unlimited")

// Unlimited represents an unlimited Limit instance.
const Unlimited = Limit(0)

// Limit represents the maximum value up to which an operation should continue processing. It may be used for defining
// the maximum number of results within a repository finder method or if applicable a template operation.
type Limit uint

// LimitOf creates a new Limit for the given max value.
func LimitOf(max uint) Limit {
	return Limit(max)
}

// Max returns the maximum number of potential results. It returns ErrUnlimited if Limit IsUnlimited.
func (l Limit) Max() (uint, error) {
	switch l {
	case Unlimited:
		return 0, ErrUnlimited
	default:
		return uint(l), nil
	}
}

// IsLimited check if the Limit is not Unlimited.
func (l Limit) IsLimited() bool {
	return l > 0
}

// IsUnlimited returns true if Limit is Unlimited.
func (l Limit) IsUnlimited() bool {
	return !l.IsLimited()
}
