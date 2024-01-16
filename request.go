package page

import (
	"errors"
	"fmt"
	"github.com/maurofran/page/sort"
	"net/http"
	"strconv"
	"strings"
)

var ErrUnpaged = errors.New("request is unpaged")

const pageParam = "page"
const sizeParam = "size"
const sortParam = "sort"

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

type parseOptions struct {
	pageParam   string
	sizeParam   string
	sortParam   string
	defaultPage uint
	defaultSize uint
	defaultSort *sort.Sort
}

// ParseOption is the type for function used to customize the RequestFrom behavior
type ParseOption func(*parseOptions) error

// WithPageParam is the option used to provide the page param.
func WithPageParam(param string) ParseOption {
	return func(options *parseOptions) error {
		if param = strings.TrimSpace(param); param != "" {
			options.pageParam = param
			return nil
		}
		return errors.New("invalid page param")
	}
}

// WithSizeParam is the option used to provide the size param.
func WithSizeParam(param string) ParseOption {
	return func(options *parseOptions) error {
		if param = strings.TrimSpace(param); param != "" {
			options.sizeParam = param
			return nil
		}
		return errors.New("invalid size param")
	}
}

// WithSortParam is the option used to provide the sort param.
func WithSortParam(param string) ParseOption {
	return func(options *parseOptions) error {
		if param = strings.TrimSpace(param); param != "" {
			options.sortParam = param
			return nil
		}
		return errors.New("invalid sort param")
	}
}

// WithDefaultPage is the option used to provide the default page number.
func WithDefaultPage(page uint) ParseOption {
	return func(options *parseOptions) error {
		options.defaultPage = page
		return nil
	}
}

// WithDefaultSize is the option used to provide the default size.
func WithDefaultSize(size uint) ParseOption {
	return func(options *parseOptions) error {
		options.defaultSize = size
		return nil
	}
}

// WithDefaultSort is the option used to provide the default sort order.
func WithDefaultSort(sorts ...string) ParseOption {
	return func(options *parseOptions) error {
		var err error
		options.defaultSort, err = sort.Parse(sorts...)
		return err
	}
}

// RequestFrom parses an http.Request query into a Request.
func RequestFrom(httpReq *http.Request, options ...ParseOption) (Request, error) {
	opts := &parseOptions{
		pageParam:   "page",
		sizeParam:   "size",
		sortParam:   "sort",
		defaultPage: 0,
		defaultSize: 10,
		defaultSort: nil,
	}
	for _, option := range options {
		if err := option(opts); err != nil {
			return nil, err
		}
	}
	query := httpReq.URL.Query()
	pageNumber := opts.defaultPage
	pageSize := opts.defaultSize
	sortClause := opts.defaultSort
	if str := query.Get(pageParam); str != "" {
		value, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, err
		}
		pageNumber = uint(value)
	}
	if str := query.Get(sizeParam); str != "" {
		value, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, err
		}
		pageSize = uint(value)
	}
	if sorts, ok := query[sortParam]; ok && len(sorts) > 0 {
		value, err := sort.Parse(sorts...)
		if err != nil {
			return nil, err
		}
		sortClause = value
	}
	return RequestOf(pageNumber, pageSize, sortClause), nil
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
