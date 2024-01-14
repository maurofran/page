package page

import (
	"errors"
	"fmt"
	"github.com/maurofran/page/sort"
)

var ErrUnpaged = errors.New("request is unpaged")

// Unpaged returns a Request instance representing no pagination setup.
func Unpaged(sorts ...*sort.Sort) Request {
	if len(sorts) == 0 {
		return unpaged{sort: sort.Unsorted()}
	}
	return unpaged{sort: sorts[0]}
}

// RequestOfSize creates a new Request for the first page (page number 0) given size.
func RequestOfSize(size uint) Request {
	return UnsortedRequestOf(0, size)
}

// UnsortedRequestOf creates a new Request for given page and size.
func UnsortedRequestOf(page, size uint) Request {
	return RequestOf(page, size, sort.Unsorted())
}

// RequestOf creates a new Request for given page and size and sort.
func RequestOf(page, size uint, sort *sort.Sort) Request {
	return request{
		page: page,
		size: size,
		sort: sort,
	}
}

// Request is the interface for pagination information.
type Request interface {
	// IsPaged returns whether the current Request contains pagination information.
	IsPaged() bool
	// IsUnpaged returns whether the current Request does not contain pagination information.
	IsUnpaged() bool
	// PageNumber returns the page to be returned. Returns ErrUnpaged if the object IsUnpaged.
	PageNumber() (uint, error)
	// PageSize returns the number of items to be returned. Returns ErrUnpaged if the object IsUnpaged.
	PageSize() (uint, error)
	// Offset returns the offset to be taken according to the underlying page and page size. Returns ErrUnpaged if the
	// object IsUnpaged.
	Offset() (uint, error)
	// Sort returns the sorting parameters.
	Sort() *sort.Sort
	// SortOr returns the current Sort or the given one if the current one is unsorted.
	SortOr(sort *sort.Sort) *sort.Sort
	// Next returns the Request requesting the next Page.
	Next() Request
	// PreviousOrFirst returns the previous Request or the first Request if the current one already is the first one.
	PreviousOrFirst() Request
	// First returns the Request requesting the first page.
	First() Request
	// WithPage returns a new Request with pageNumber applied. It returns ErrUnpaged if the object IsUnpaged and the
	// pageNumber is not zero.
	WithPage(pageNumber uint) (Request, error)
	// HasPrevious returns whether there's a previous Request we can access from the current one. Will return false in
	// case the current Request already refers to the first page.
	HasPrevious() bool
	// ToLimit returns a Limit from this Request if IsPaged or Unlimited otherwise.
	ToLimit() Limit
	// ToScrollPosition returns an OffsetScrollPosition from this Request if IsPaged. It returns ErrUnpaged otherwise.
	ToScrollPosition() (OffsetScrollPosition, error)
}

type unpaged struct {
	sort *sort.Sort
}

func (u unpaged) IsPaged() bool {
	return false
}

func (u unpaged) IsUnpaged() bool {
	return true
}

func (u unpaged) PageNumber() (uint, error) {
	return 0, ErrUnpaged
}

func (u unpaged) PageSize() (uint, error) {
	return 0, ErrUnpaged
}

func (u unpaged) Offset() (uint, error) {
	return 0, ErrUnpaged
}

func (u unpaged) Sort() *sort.Sort {
	return u.sort
}

func (u unpaged) SortOr(sort *sort.Sort) *sort.Sort {
	if s := u.sort; s != nil && s.IsSorted() {
		return s
	}
	return sort
}

func (u unpaged) Next() Request {
	return u
}

func (u unpaged) PreviousOrFirst() Request {
	return u
}

func (u unpaged) First() Request {
	return u
}

func (u unpaged) WithPage(pageNumber uint) (Request, error) {
	if pageNumber == 0 {
		return u, nil
	}
	return unpaged{}, ErrUnpaged
}

func (u unpaged) HasPrevious() bool {
	return false
}

func (u unpaged) ToLimit() Limit {
	return Unlimited
}

func (u unpaged) ToScrollPosition() (OffsetScrollPosition, error) {
	return InitialOffsetScrollPosition, ErrUnpaged
}

type request struct {
	page uint
	size uint
	sort *sort.Sort
}

func (r request) IsPaged() bool {
	return true
}

func (r request) IsUnpaged() bool {
	return false
}

func (r request) PageNumber() (uint, error) {
	return r.page, nil
}

func (r request) PageSize() (uint, error) {
	return r.size, nil
}

func (r request) Offset() (uint, error) {
	return r.page * r.size, nil
}

func (r request) Sort() *sort.Sort {
	return r.sort
}

func (r request) SortOr(sort *sort.Sort) *sort.Sort {
	if r.sort.IsUnsorted() {
		return sort
	}
	return r.sort
}

func (r request) Next() Request {
	return request{
		page: r.page + 1,
		size: r.size,
		sort: r.sort,
	}
}

func (r request) Previous() Request {
	return request{
		page: r.page - 1,
		size: r.size,
		sort: r.sort,
	}
}

func (r request) PreviousOrFirst() Request {
	if r.HasPrevious() {
		return r.Previous()
	}
	return r.First()
}

func (r request) First() Request {
	return request{
		page: 0,
		size: r.size,
		sort: r.sort,
	}
}

func (r request) WithPage(pageNumber uint) (Request, error) {
	return request{
		page: pageNumber,
		size: r.size,
		sort: r.sort,
	}, nil
}

func (r request) HasPrevious() bool {
	return r.page > 0
}

func (r request) ToLimit() Limit {
	return LimitOf(r.size)
}

func (r request) ToScrollPosition() (OffsetScrollPosition, error) {
	offset, _ := r.Offset()
	return OffsetScrollPositionOf(offset), nil
}

func (r request) String() string {
	return fmt.Sprintf(
		"Page request [number: %d, size: %d, sort: %s]",
		r.page,
		r.size,
		r.sort,
	)
}
