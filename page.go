package page

import (
	"encoding/json"
	"fmt"
	"github.com/maurofran/page/sort"
	"math"
	"reflect"
)

// Empty returns an empty page with provided (optional) page Request.
// If more than one Request is provided, the first one will be used and the others ignored.
func Empty[T any](request ...Request) Page[T] {
	req := Unpaged()
	if len(request) > 0 {
		req = request[0]
	}
	page := &page[T]{
		chunk: chunk[T]{
			request: req,
		},
		totalElements: 0,
	}
	page.chunk.hasPrev = page.hasPrevious
	page.chunk.hasNext = page.hasNext
	return page
}

// Page is a sublist of a list of objects. It allows gain information about the position of it in the containing
// entire list.
type Page[T any] interface {
	Slice[T]

	// TotalPages returns the total number of pages.
	TotalPages() uint
	// TotalElements returns the total amount of elements.
	TotalElements() uint
}

type page[T any] struct {
	chunk[T]
	totalElements uint
}

// New creates new Page.
func New[T any](content []T, request Request, totalElements uint) Page[T] {
	if request != nil && len(content) > 0 && request.IsPaged() {
		offset, _ := request.Offset()
		size, _ := request.PageSize()
		if (offset + size) > totalElements {
			totalElements = offset + size
		}
	}
	page := &page[T]{
		chunk: chunk[T]{
			content: content,
			request: request,
		},
		totalElements: totalElements,
	}
	page.chunk.hasPrev = page.hasPrevious
	page.chunk.hasNext = page.hasNext
	return page
}

// FromSlice creates an Unpaged Page with supplied content.
func FromSlice[T any](content []T) Page[T] {
	return New(content, Unpaged(), uint(len(content)))
}

func (p *page[T]) TotalPages() uint {
	if size := p.Size(); size > 0 {
		return uint(math.Ceil(float64(p.totalElements) / float64(size)))
	}
	return 1
}

func (p *page[T]) TotalElements() uint {
	return p.totalElements
}

func (p *page[T]) hasPrevious() bool {
	return p.Number() > 0
}

func (p *page[T]) hasNext() bool {
	return (p.Number() + 1) < p.TotalPages()
}

type jsonPage[T any] struct {
	Content       []T        `json:"content"`
	Number        uint       `json:"number"`
	Size          uint       `json:"size"`
	Sort          *sort.Sort `json:"sort"`
	TotalElements uint       `json:"totalElements"`
	TotalPages    uint       `json:"totalPages"`
}

// MarshalJSON implements the json.Marshaler interface.
func (p *page[T]) MarshalJSON() ([]byte, error) {
	payload := jsonPage[T]{
		Content:       p.Content(),
		Number:        p.Number(),
		Size:          p.Size(),
		Sort:          p.Sort(),
		TotalElements: p.TotalElements(),
		TotalPages:    p.TotalPages(),
	}
	return json.Marshal(payload)
}

func (p *page[T]) UnmarshalJSON(data []byte) error {
	var payload jsonPage[T]
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	p.request = RequestOf(payload.Number, payload.Size, payload.Sort)
	p.content = payload.Content
	p.totalElements = payload.TotalElements
	return nil
}

func (p *page[T]) String() string {
	contentType := "UNKNOWN"
	if content := p.content; len(content) > 0 {
		contentType = reflect.TypeOf(content[0]).String()
	}
	return fmt.Sprintf(
		"Page %d of %d containing %s instances",
		p.Number()+1,
		p.TotalPages(),
		contentType,
	)
}
