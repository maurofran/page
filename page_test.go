package page_test

import (
	"github.com/maurofran/page"
	"github.com/maurofran/page/sort"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEmpty(t *testing.T) {
	t.Run("unpaged", func(t *testing.T) {
		fixture := page.Empty[int]()
		assert.Len(t, fixture.Content(), 0)
		assert.Equal(t, uint(1), fixture.TotalPages())
		assert.Equal(t, uint(0), fixture.TotalElements())
		assert.Equal(t, page.Unpaged(), fixture.Request())
	})

	t.Run("paged", func(t *testing.T) {
		fixture := page.Empty[int](page.RequestOf(1, 10))
		assert.Len(t, fixture.Content(), 0)
		assert.Equal(t, uint(0), fixture.TotalPages())
		assert.Equal(t, uint(0), fixture.TotalElements())
		assert.Equal(t, page.RequestOf(1, 10), fixture.Request())
	})
}

func TestNew(t *testing.T) {
	t.Run("paged", func(t *testing.T) {
		content := []int{1, 2, 3, 4, 5}
		request := page.RequestOf(0, 5)
		fixture := page.New(content, request, 56)
		assert.Equal(t, content, fixture.Content())
		assert.Equal(t, request, fixture.Request())
		assert.Equal(t, uint(12), fixture.TotalPages())
		assert.Equal(t, uint(56), fixture.TotalElements())
	})
	t.Run("total elements", func(t *testing.T) {
		content := []int{1, 2, 3, 4, 5}
		request := page.RequestOf(0, 5)
		fixture := page.New(content, request, 0)
		assert.Equal(t, content, fixture.Content())
		assert.Equal(t, request, fixture.Request())
		assert.Equal(t, uint(1), fixture.TotalPages())
		assert.Equal(t, uint(5), fixture.TotalElements())
	})
}

func TestFromSlice(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		fixture := page.FromSlice[int](nil)
		assert.Len(t, fixture.Content(), 0)
		assert.Equal(t, uint(1), fixture.TotalPages())
		assert.Equal(t, uint(0), fixture.TotalElements())
		assert.Equal(t, page.Unpaged(), fixture.Request())
	})

	t.Run("content", func(t *testing.T) {
		fixture := page.FromSlice[int]([]int{1, 2, 3, 4, 5})
		assert.Len(t, fixture.Content(), 5)
		assert.Equal(t, uint(1), fixture.TotalPages())
		assert.Equal(t, uint(5), fixture.TotalElements())
		assert.Equal(t, page.Unpaged(), fixture.Request())
	})
}

func TestPage_Number(t *testing.T) {
	t.Run("unpaged", func(t *testing.T) {
		fixture := page.New([]int{1, 2, 3, 4, 5}, page.Unpaged(), 56)
		assert.Equal(t, uint(0), fixture.Number())
	})

	t.Run("paged", func(t *testing.T) {
		fixture := page.New([]int{1, 2, 3, 4, 5}, page.RequestOf(4, 5), 56)
		assert.Equal(t, uint(4), fixture.Number())
	})
}

func TestPage_Size(t *testing.T) {
	t.Run("unpaged", func(t *testing.T) {
		fixture := page.New([]int{1, 2, 3, 4, 5}, page.Unpaged(), 56)
		assert.Equal(t, uint(5), fixture.Size())
	})

	t.Run("paged", func(t *testing.T) {
		fixture := page.New([]int{1, 2, 3, 4, 5}, page.RequestOf(4, 10), 56)
		assert.Equal(t, uint(10), fixture.Size())
	})
}

func TestPage_NumberOfElements(t *testing.T) {
	fixture := page.FromSlice[int]([]int{1, 2, 3})
	assert.Equal(t, uint(3), fixture.NumberOfElements())
}

func TestPage_Content(t *testing.T) {
	fixture := page.FromSlice[int]([]int{1, 2, 3})
	assert.Equal(t, []int{1, 2, 3}, fixture.Content())
}

func TestPage_HasContent(t *testing.T) {
	t.Run("no content", func(t *testing.T) {
		fixture := page.Empty[int]()
		assert.False(t, fixture.HasContent())
	})

	t.Run("has content", func(t *testing.T) {
		fixture := page.FromSlice[int]([]int{1, 2, 3})
		assert.True(t, fixture.HasContent())
	})
}

func TestPage_Sort(t *testing.T) {
	t.Run("unpaged", func(t *testing.T) {
		fixture := page.Empty[int]()
		assert.Equal(t, sort.Unsorted(), fixture.Sort())
	})

	t.Run("paged", func(t *testing.T) {
		request := parseRequest(t, "page=4&size=10&sort=id,desc")
		fixture := page.New([]int{1, 2, 3, 4, 5}, request, 56)
		expected, err := sort.Parse("id,desc")
		require.NoError(t, err)
		assert.Equal(t, expected, fixture.Sort())
	})
}

func TestPage_IsFirst(t *testing.T) {
	t.Run("unpaged", func(t *testing.T) {
		fixture := page.Empty[int]()
		assert.True(t, fixture.IsFirst())
	})

	t.Run("paged first", func(t *testing.T) {
		fixture := page.New([]int{1, 2, 3, 4, 5}, page.RequestOf(0, 5), 56)
		assert.True(t, fixture.IsFirst())
	})

	t.Run("paged other", func(t *testing.T) {
		fixture := page.New([]int{1, 2, 3, 4, 5}, page.RequestOf(1, 5), 56)
		assert.False(t, fixture.IsFirst())
	})
}

func TestPage_IsLast(t *testing.T) {
	t.Run("unpaged", func(t *testing.T) {
		fixture := page.Empty[int]()
		assert.True(t, fixture.IsLast())
	})

	t.Run("paged last", func(t *testing.T) {
		fixture := page.New([]int{1, 2, 3, 4, 5}, page.RequestOf(11, 5), 56)
		assert.True(t, fixture.IsLast())
	})

	t.Run("paged other", func(t *testing.T) {
		fixture := page.New([]int{1, 2, 3, 4, 5}, page.RequestOf(1, 5), 56)
		assert.False(t, fixture.IsLast())
	})
}

func TestPage_NextRequest(t *testing.T) {
	t.Run("unpaged", func(t *testing.T) {
		fixture := page.Empty[int]()
		assert.Equal(t, page.Unpaged(), fixture.NextRequest())
	})

	t.Run("paged", func(t *testing.T) {
		fixture := page.New([]int{1, 2, 3, 4, 5}, page.RequestOf(4, 5), 56)
		assert.Equal(t, page.RequestOf(5, 5), fixture.NextRequest())
	})
}

func TestPage_PreviousRequest(t *testing.T) {
	t.Run("unpaged", func(t *testing.T) {
		fixture := page.Empty[int]()
		assert.Equal(t, page.Unpaged(), fixture.PreviousRequest())
	})

	t.Run("paged", func(t *testing.T) {
		fixture := page.New([]int{1, 2, 3, 4, 5}, page.RequestOf(4, 5), 56)
		assert.Equal(t, page.RequestOf(3, 5), fixture.PreviousRequest())
	})
}

func TestPage_NextOrLastRequest(t *testing.T) {
	t.Run("unpaged", func(t *testing.T) {
		fixture := page.Empty[int]()
		assert.Equal(t, page.Unpaged(), fixture.NextOrLastRequest())
	})

	t.Run("paged", func(t *testing.T) {
		fixture := page.New([]int{1, 2, 3, 4, 5}, page.RequestOf(4, 5), 56)
		assert.Equal(t, page.RequestOf(5, 5), fixture.NextOrLastRequest())
	})

	t.Run("last", func(t *testing.T) {
		fixture := page.New([]int{1, 2, 3, 4, 5}, page.RequestOf(4, 5), 20)
		assert.Equal(t, page.RequestOf(4, 5), fixture.NextOrLastRequest())
	})
}

func TestPage_PreviousOrFirstRequest(t *testing.T) {
	t.Run("unpaged", func(t *testing.T) {
		fixture := page.Empty[int]()
		assert.Equal(t, page.Unpaged(), fixture.PreviousOrFirstRequest())
	})

	t.Run("paged", func(t *testing.T) {
		fixture := page.New([]int{1, 2, 3, 4, 5}, page.RequestOf(4, 5), 56)
		assert.Equal(t, page.RequestOf(3, 5), fixture.PreviousOrFirstRequest())
	})

	t.Run("last", func(t *testing.T) {
		fixture := page.New([]int{1, 2, 3, 4, 5}, page.RequestOf(0, 5), 20)
		assert.Equal(t, page.RequestOf(0, 5), fixture.PreviousOrFirstRequest())
	})
}
