package mgo

import (
	"errors"
	"math"
	"testing"

	"github.com/jucardi/go-mongodb-lib/pages"
	"github.com/jucardi/go-mongodb-lib/testutils"
	"github.com/stretchr/testify/assert"
)

func TestQuery_Page_NoSort(t *testing.T) {
	q := MockQuery(t)

	q.When("Page", makePageHandler(q))
	q.BulkTimes([]string{"Sort", "Skip", "Limit"}, []int{0, 0, 0})

	page := &pages.Page{Page: 1, Size: 10}
	q.Page(page)

	q.BulkTimes([]string{"Sort", "Skip", "Limit"}, []int{0, 1, 1})
}

func TestQuery_Page_Sort(t *testing.T) {
	q := MockQuery(t)

	q.When("Page", makePageHandler(q))
	q.BulkTimes([]string{"Sort", "Skip", "Limit"}, []int{0, 0, 0})

	page := &pages.Page{Page: 1, Size: 10, Sort: []string{"name", "other"}}
	q.Page(page)

	q.BulkTimes([]string{"Sort", "Skip", "Limit"}, []int{1, 1, 1})
}

func TestQuery_Page_NoPage(t *testing.T) {
	q := MockQuery(t)

	q.When("Page", makePageHandler(q))
	q.BulkTimes([]string{"Sort", "Skip", "Limit"}, []int{0, 0, 0})

	q.Page()

	q.BulkTimes([]string{"Sort", "Skip", "Limit"}, []int{0, 0, 0})
}

type testObj struct {
	A string
	B int
}

const (
	testString = "answer to life"
	testInt    = 42
)

func TestQuery_WrapPage_Success_NoSort(t *testing.T) {
	page := &pages.Page{Page: 1, Size: 10}
	count := 100
	totalPages := int(math.Ceil(float64(count) / float64(page.Size)))
	q := MockQuery(t)

	q.BulkTimes([]string{"Sort", "Skip", "Limit", "Count", "Page"}, []int{0, 0, 0, 0, 0})
	q.WhenReturn("Count", count, nil)
	q.When("Page", makePageHandler(q))
	q.When("WrapPage", makeWrapPageHandler(q))
	q.When("All", makeAllHandler(page.Size))

	var ret []*testObj

	paginated, err := q.WrapPage(&ret, page)
	q.BulkTimes([]string{"Sort", "Skip", "Limit", "Count", "Page"}, []int{0, 1, 1, 1, 1})
	assert.Nil(t, err)
	assert.NotNil(t, paginated)

	assert.Equal(t, page.Size, paginated.ItemsCount)
	assert.Equal(t, count, paginated.TotalCount)
	assert.Equal(t, totalPages, paginated.TotalPages)
	assert.Equal(t, page.Page, paginated.Page)
	assert.Equal(t, page.Size, paginated.Size)
	assert.Len(t, paginated.Items, page.Size)
	assert.Len(t, ret, page.Size)

	result := ret[0]
	assert.Equal(t, testString, result.A)
	assert.Equal(t, testInt, result.B)
	println(result.A, result.B)
}

func TestQuery_WrapPage_Success_NoPage(t *testing.T) {
	count := 100
	q := MockQuery(t)

	q.BulkTimes([]string{"Sort", "Skip", "Limit", "Count", "Page"}, []int{0, 0, 0, 0, 0})
	q.WhenReturn("Count", count, nil)
	q.When("Page", makePageHandler(q))
	q.When("WrapPage", makeWrapPageHandler(q))
	q.When("All", makeAllHandler(count))

	var ret []*testObj

	paginated, err := q.WrapPage(&ret)
	q.BulkTimes([]string{"Sort", "Skip", "Limit", "Count", "Page"}, []int{0, 0, 0, 0, 0})
	assert.Nil(t, err)
	assert.NotNil(t, paginated)

	assert.Equal(t, count, paginated.ItemsCount)
	assert.Equal(t, count, paginated.TotalCount)
	assert.Equal(t, 1, paginated.TotalPages)
	assert.Equal(t, 1, paginated.Page)
	assert.Equal(t, count, paginated.Size)
	assert.Len(t, paginated.Items, count)
	assert.Len(t, ret, count)

	result := ret[0]
	assert.Equal(t, testString, result.A)
	assert.Equal(t, testInt, result.B)
	println(result.A, result.B)
}

func TestQuery_WrapPage_Fail_CountError(t *testing.T) {
	page := &pages.Page{Page: 1, Size: 10}
	q := MockQuery(t)

	q.BulkTimes([]string{"Sort", "Skip", "Limit", "Count", "Page"}, []int{0, 0, 0, 0, 0})
	q.WhenReturn("Count", 0, errors.New("forgot how to count"))
	q.When("Page", makePageHandler(q))
	q.When("WrapPage", makeWrapPageHandler(q))
	q.When("All", makeAllHandler(page.Size))

	var ret []*testObj

	paginated, err := q.WrapPage(&ret, page)
	q.BulkTimes([]string{"Sort", "Skip", "Limit", "Count", "Page"}, []int{0, 0, 0, 1, 0})
	assert.NotNil(t, err)
	assert.Equal(t, "unable to obtain a count of elements, forgot how to count", err.Error())
	assert.Nil(t, paginated)
	assert.Len(t, ret, 0)
}

func TestQuery_WrapPage_Fail_AllError(t *testing.T) {
	count := 100
	q := MockQuery(t)

	q.BulkTimes([]string{"Sort", "Skip", "Limit", "Count", "Page"}, []int{0, 0, 0, 0, 0})
	q.WhenReturn("Count", count, nil)
	q.When("Page", makePageHandler(q))
	q.When("WrapPage", makeWrapPageHandler(q))
	q.WhenReturn("All", errors.New("i just broke"))

	var ret []*testObj

	paginated, err := q.WrapPage(&ret)
	q.BulkTimes([]string{"Sort", "Skip", "Limit", "Count", "Page"}, []int{0, 0, 0, 0, 0})
	assert.NotNil(t, err)
	assert.Equal(t, "failed to unmarshal page, i just broke", err.Error())
	assert.Nil(t, paginated)
	assert.Len(t, ret, 0)
}

func makePageHandler(q IQuery) WhenHandler {
	return func(t *testing.T, args ...interface{}) []interface{} {
		var page *pages.Page
		if len(args) > 0 && args[0] != nil {
			page = args[0].(*pages.Page)
		}
		return testutils.MakeReturn(pageHandler(q, page))
	}
}

func makeWrapPageHandler(q IQuery) WhenHandler {
	return func(t *testing.T, args ...interface{}) []interface{} {
		result := args[0]
		var page *pages.Page
		if len(args) > 1 && args[1] != nil {
			page = args[1].(*pages.Page)
		}
		return testutils.MakeReturn(wrapPageHandler(q, result, page))
	}
}

func makeAllHandler(size int) WhenHandler {
	return func(t *testing.T, args ...interface{}) []interface{} {
		result := args[0]
		list := result.(*[]*testObj)

		for i := 0; i < size; i++ {
			*list = append(*list, &testObj{
				A: testString,
				B: testInt,
			})
		}

		return testutils.MakeReturn(nil)
	}
}
