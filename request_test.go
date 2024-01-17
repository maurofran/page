package page_test

import (
	"fmt"
	"github.com/maurofran/page"
	"github.com/maurofran/page/sort"
	"github.com/maurofran/page/sort/order/direction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

func TestUnpaged(t *testing.T) {
	t.Run("unsorted", func(t *testing.T) {
		fixture := page.Unpaged()
		assert.True(t, fixture.IsUnpaged())
		assert.Equal(t, sort.Unsorted(), fixture.Sort())
	})

	t.Run("sorted", func(t *testing.T) {
		fixture := page.Unpaged(sort.By(direction.Asc, "name"))
		assert.True(t, fixture.IsUnpaged())
		assert.Equal(t, sort.By(direction.Asc, "name"), fixture.Sort())
	})
}

func TestRequestOfSize(t *testing.T) {
	fixture := page.RequestOfSize(10)
	assert.False(t, fixture.IsUnpaged())
	pageNumber, err := fixture.PageNumber()
	assert.NoError(t, err)
	assert.Equal(t, uint(0), pageNumber)
	pageSize, err := fixture.PageSize()
	assert.NoError(t, err)
	assert.Equal(t, uint(10), pageSize)
	sorting := fixture.Sort()
	assert.True(t, sorting.IsUnsorted())
}

func TestRequestOf(t *testing.T) {
	fixture := page.RequestOf(2, 20, sort.By(direction.Desc, "name"))
	assert.False(t, fixture.IsUnpaged())
	pageNumber, err := fixture.PageNumber()
	assert.NoError(t, err)
	assert.Equal(t, uint(2), pageNumber)
	pageSize, err := fixture.PageSize()
	assert.NoError(t, err)
	assert.Equal(t, uint(20), pageSize)
	sorting := fixture.Sort()
	assert.False(t, sorting.IsUnsorted())
	assert.Equal(t, sort.By(direction.Desc, "name"), sorting)
}

func TestRequestFrom(t *testing.T) {
	tests := map[string]struct {
		query string
		page  uint
		size  uint
		sort  *sort.Sort
		err   error
	}{
		"empty":       {"", 0, 10, nil, nil},
		"invalidPage": {"page=foo", 0, 10, nil, strconv.ErrSyntax},
		"invalidSize": {"page=2&size=foo", 2, 10, nil, strconv.ErrSyntax},
		"invalidSort": {"sort=,invalid", 0, 10, nil, direction.ErrInvalid},
		"allParams":   {"page=2&size=20&sort=name,desc", 2, 20, sort.By(direction.Desc, "name"), nil},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			request := &http.Request{URL: &url.URL{RawQuery: test.query}}
			fixture, err := page.RequestFrom(request)
			assert.ErrorIs(t, err, test.err)
			if err == nil {
				pageNumber, err := fixture.PageNumber()
				assert.NoError(t, err)
				assert.Equal(t, test.page, pageNumber)
				pageSize, err := fixture.PageSize()
				assert.NoError(t, err)
				assert.Equal(t, test.size, pageSize)
				sorting := fixture.Sort()
				assert.Equal(t, sorting, test.sort)
			}
		})
	}
}

func TestRequestFromCustomOptions(t *testing.T) {
	tests := map[string]struct {
		query string
		page  uint
		size  uint
		sort  *sort.Sort
		err   error
	}{
		"empty":     {"", 1, 50, sort.By(direction.Asc, "name"), nil},
		"allParams": {"thePage=2&theSize=20&theSort=name,desc", 2, 20, sort.By(direction.Desc, "name"), nil},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			request := &http.Request{URL: &url.URL{RawQuery: test.query}}
			fixture, err := page.RequestFrom(
				request,
				page.WithDefaultPage(1),
				page.WithDefaultSize(50),
				page.WithDefaultSort("name,asc"),
				page.WithPageParam("thePage"),
				page.WithSizeParam("theSize"),
				page.WithSortParam("theSort"),
			)
			assert.ErrorIs(t, err, test.err)
			if err == nil {
				pageNumber, err := fixture.PageNumber()
				assert.NoError(t, err)
				assert.Equal(t, test.page, pageNumber)
				pageSize, err := fixture.PageSize()
				assert.NoError(t, err)
				assert.Equal(t, test.size, pageSize)
				sorting := fixture.Sort()
				assert.Equal(t, sorting, test.sort)
			}
		})
	}
}

func TestUnpaged_IsPaged(t *testing.T) {
	fixture := page.Unpaged()
	assert.False(t, fixture.IsPaged())
}

func TestUnpaged_IsUnpaged(t *testing.T) {
	fixture := page.Unpaged()
	assert.True(t, fixture.IsUnpaged())
}

func TestUnpaged_PageNumber(t *testing.T) {
	fixture := page.Unpaged()
	_, err := fixture.PageNumber()
	assert.ErrorIs(t, err, page.ErrUnpaged)
}

func TestUnpaged_PageSize(t *testing.T) {
	fixture := page.Unpaged()
	_, err := fixture.PageSize()
	assert.ErrorIs(t, err, page.ErrUnpaged)
}

func TestUnpaged_Offset(t *testing.T) {
	fixture := page.Unpaged()
	_, err := fixture.Offset()
	assert.ErrorIs(t, err, page.ErrUnpaged)
}

func TestUnpaged_Sort(t *testing.T) {
	fixture := page.Unpaged(sort.By(direction.Asc, "name"))
	sorting := fixture.Sort()
	assert.Equal(t, sorting, sort.By(direction.Asc, "name"))
}

func TestUnpaged_SortOr(t *testing.T) {
	t.Run("unsorted", func(t *testing.T) {
		fixture := page.Unpaged()
		sorting := fixture.SortOr(sort.By(direction.Asc, "name"))
		assert.Equal(t, sorting, sort.By(direction.Asc, "name"))
	})

	t.Run("sorted", func(t *testing.T) {
		fixture := page.Unpaged(sort.By(direction.Desc, "age"))
		sorting := fixture.SortOr(sort.By(direction.Asc, "name"))
		assert.Equal(t, sorting, sort.By(direction.Desc, "age"))

	})
}

func TestUnpaged_Next(t *testing.T) {
	fixture := page.Unpaged()
	assert.Equal(t, fixture, fixture.Next())
}

func TestUnpaged_PreviousOrFirst(t *testing.T) {
	fixture := page.Unpaged()
	assert.Equal(t, fixture, fixture.PreviousOrFirst())
}

func TestUnpaged_First(t *testing.T) {
	fixture := page.Unpaged()
	assert.Equal(t, fixture, fixture.First())
}

func TestUnpaged_WithPage(t *testing.T) {
	fixture := page.Unpaged()
	t.Run("page_0", func(t *testing.T) {
		actual, err := fixture.WithPage(0)
		assert.NoError(t, err)
		assert.Equal(t, fixture, actual)
	})
	t.Run("page_1", func(t *testing.T) {
		_, err := fixture.WithPage(1)
		assert.ErrorIs(t, err, page.ErrUnpaged)
	})
}

func TestUnpaged_HasPrevious(t *testing.T) {
	fixture := page.Unpaged()
	assert.False(t, fixture.HasPrevious())
}

func TestUnpaged_ToLimit(t *testing.T) {
	fixture := page.Unpaged()
	limit := fixture.ToLimit()
	assert.Equal(t, page.Unlimited, limit)
}

func TestUnpaged_ToScrollPosition(t *testing.T) {
	fixture := page.Unpaged()
	_, err := fixture.ToScrollPosition()
	assert.ErrorIs(t, err, page.ErrUnpaged)
}

func parseRequest(t *testing.T, query string) page.Request {
	t.Helper()
	httpRequest := &http.Request{URL: &url.URL{RawQuery: query}}
	request, err := page.RequestFrom(httpRequest)
	require.NoError(t, err)
	return request
}

func TestRequest_IsPaged(t *testing.T) {
	request := parseRequest(t, "page=0&size=10")
	assert.True(t, request.IsPaged())
}

func TestRequest_IsUnpaged(t *testing.T) {
	request := parseRequest(t, "page=0&size=10")
	assert.False(t, request.IsUnpaged())
}

func TestRequest_PageNumber(t *testing.T) {
	request := parseRequest(t, "page=0&size=10")
	pageNumber, err := request.PageNumber()
	require.NoError(t, err)
	assert.Equal(t, uint(0), pageNumber)
}

func TestRequest_PageSize(t *testing.T) {
	request := parseRequest(t, "page=0&size=10")
	pageSize, err := request.PageSize()
	require.NoError(t, err)
	assert.Equal(t, uint(10), pageSize)
}

func TestRequest_Offset(t *testing.T) {
	request := parseRequest(t, "page=2&size=10")
	offset, err := request.Offset()
	require.NoError(t, err)
	assert.Equal(t, uint(20), offset)
}

func TestRequest_Sort(t *testing.T) {
	request := parseRequest(t, "page=0&size=10&sort=name,asc")
	sorting := request.Sort()
	assert.Equal(t, sorting, sort.By(direction.Asc, "name"))
}

func TestRequest_SortOr(t *testing.T) {
	t.Run("unsorted", func(t *testing.T) {
		request := parseRequest(t, "page=0&size=10")
		sorting := request.SortOr(sort.By(direction.Asc, "name"))
		assert.Equal(t, sorting, sort.By(direction.Asc, "name"))
	})

	t.Run("sorted", func(t *testing.T) {
		request := parseRequest(t, "page=0&size=10&sort=age,desc")
		sorting := request.SortOr(sort.By(direction.Asc, "name"))
		assert.Equal(t, sorting, sort.By(direction.Desc, "age"))

	})
}

func TestRequest_Next(t *testing.T) {
	request := parseRequest(t, "page=0&size=10&sort=name,asc")
	actual := request.Next()
	pageNumber, _ := actual.PageNumber()
	assert.Equal(t, uint(1), pageNumber)
	pageSize, _ := actual.PageSize()
	assert.Equal(t, uint(10), pageSize)
	sorting := actual.Sort()
	assert.Equal(t, sorting, sort.By(direction.Asc, "name"))
}

func TestRequest_Previous(t *testing.T) {
	request := parseRequest(t, "page=3&size=10&sort=name,asc")
	actual := request.Next()
	pageNumber, _ := actual.PageNumber()
	assert.Equal(t, uint(4), pageNumber)
	pageSize, _ := actual.PageSize()
	assert.Equal(t, uint(10), pageSize)
	sorting := actual.Sort()
	assert.Equal(t, sorting, sort.By(direction.Asc, "name"))
}

func TestRequest_PreviousOrFirst(t *testing.T) {
	t.Run("has_previous", func(t *testing.T) {
		request := parseRequest(t, "page=3&size=10&sort=name,asc")
		actual := request.PreviousOrFirst()
		pageNumber, _ := actual.PageNumber()
		assert.Equal(t, uint(2), pageNumber)
		pageSize, _ := actual.PageSize()
		assert.Equal(t, uint(10), pageSize)
		sorting := actual.Sort()
		assert.Equal(t, sorting, sort.By(direction.Asc, "name"))
	})

	t.Run("first", func(t *testing.T) {
		request := parseRequest(t, "page=0&size=10&sort=name,asc")
		actual := request.PreviousOrFirst()
		pageNumber, _ := actual.PageNumber()
		assert.Equal(t, uint(0), pageNumber)
		pageSize, _ := actual.PageSize()
		assert.Equal(t, uint(10), pageSize)
		sorting := actual.Sort()
		assert.Equal(t, sorting, sort.By(direction.Asc, "name"))
	})
}

func TestRequest_First(t *testing.T) {
	request := parseRequest(t, "page=9&size=10&sort=name,asc")
	actual := request.First()
	pageNumber, _ := actual.PageNumber()
	assert.Equal(t, uint(0), pageNumber)
	pageSize, _ := actual.PageSize()
	assert.Equal(t, uint(10), pageSize)
	sorting := actual.Sort()
	assert.Equal(t, sorting, sort.By(direction.Asc, "name"))
}

func TestRequest_WithPage(t *testing.T) {
	request := parseRequest(t, "page=9&size=10&sort=name,asc")
	actual, err := request.WithPage(6)
	assert.NoError(t, err)
	pageNumber, _ := actual.PageNumber()
	assert.Equal(t, uint(6), pageNumber)
	pageSize, _ := actual.PageSize()
	assert.Equal(t, uint(10), pageSize)
	sorting := actual.Sort()
	assert.Equal(t, sorting, sort.By(direction.Asc, "name"))
}

func TestRequest_HasPrevious(t *testing.T) {
	t.Run("has_previous", func(t *testing.T) {
		request := parseRequest(t, "page=9&size=10&sort=name,asc")
		assert.True(t, request.HasPrevious())
	})

	t.Run("no_previous", func(t *testing.T) {
		request := parseRequest(t, "page=0&size=10&sort=name,asc")
		assert.False(t, request.HasPrevious())
	})
}

func TestRequest_ToLimit(t *testing.T) {
	request := parseRequest(t, "page=9&size=10&sort=name,asc")
	limit := request.ToLimit()
	assert.Equal(t, page.LimitOf(uint(10)), limit)
}

func TestRequest_ToScrollPosition(t *testing.T) {
	request := parseRequest(t, "page=9&size=10&sort=name,asc")
	scrollPosition, err := request.ToScrollPosition()
	require.NoError(t, err)
	assert.Equal(t, page.OffsetScrollPositionOf(uint(90)), scrollPosition)
}

func TestRequest_String(t *testing.T) {
	request := parseRequest(t, "page=9&size=10&sort=name,asc")
	str := fmt.Sprintf("%s", request)
	assert.Equal(t, "Page request [number: 9, size: 10, sort: name: ASC]", str)
}
