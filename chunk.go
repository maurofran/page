package page

import (
	"github.com/maurofran/page/sort"
)

type chunk[T any] struct {
	content []T
	request Request
}

func (c *chunk[T]) Number() uint {
	number, err := c.request.PageNumber()
	if err != nil {
		return 0
	}
	return number
}

func (c *chunk[T]) Size() uint {
	size, err := c.request.PageSize()
	if err != nil {
		return c.NumberOfElements()
	}
	return size
}

func (c *chunk[T]) NumberOfElements() uint {
	return uint(len(c.content))
}

func (c *chunk[T]) Content() []T {
	return c.content
}

func (c *chunk[T]) HasContent() bool {
	return len(c.content) > 0
}

func (c *chunk[T]) Sort() *sort.Sort {
	return c.request.Sort()
}

func (c *chunk[T]) IsFirst() bool {
	return !c.HasPrevious()
}

func (c *chunk[T]) IsLast() bool {
	return !c.HasNext()
}

func (c *chunk[T]) HasNext() bool {
	panic("implement me")
}

func (c *chunk[T]) HasPrevious() bool {
	panic("implement me")
}

func (c *chunk[T]) Request() Request {
	return c.request
}

func (c *chunk[T]) NextPageable() Request {
	if c.HasNext() {
		return c.request.Next()
	}
	return Unpaged(c.Sort())
}

func (c *chunk[T]) PreviousPageable() Request {
	if c.HasPrevious() {
		return c.request.PreviousOrFirst()
	}
	return Unpaged(c.Sort())
}

func (c *chunk[T]) NextOrLastPageable() Request {
	if c.HasNext() {
		return c.NextPageable()
	}
	return c.Request()
}

func (c *chunk[T]) PreviousOrFirstPageable() Request {
	if c.HasPrevious() {
		return c.PreviousPageable()
	}
	return c.Request()
}
