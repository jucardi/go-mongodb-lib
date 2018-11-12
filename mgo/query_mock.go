package mgo

import (
	"github.com/jucardi/go-mongodb-lib/pages"
)

// TODO: Handle conditionals on arguments to support more use cases when testing. Add all IQuery methods.

// All functions that return IQuery to easily initialize
var queryFuncs = []string{"Batch", "Prefetch", "Skip", "Limit", "Select", "Sort", "Count", "Explain", "Hint", "SetMaxScan", "SetMaxTime", "Snapshot", "Comment", "LogReplay", "Page"}

// QueryMock is a mock implementation of IQuery
type QueryMock struct {
	IQuery
	*MockBase
}

// MockQuery returns a new instance of IQuery for mocking purposes
//
//   {q}:   An optional instance of IQuery to handle real calls if wanted.
//
func MockQuery(q ...IQuery) *QueryMock {
	ret := &QueryMock{
		IQuery:   NewQuery(),
		MockBase: newMock(),
	}
	if len(q) > 0 {
		ret.IQuery = q[0]
	}
	return ret.init()
}

func (m *QueryMock) init() *QueryMock {
	// Sets the default return value for the methods that return IQuery. By default they'll return the mock.
	for _, v := range queryFuncs {
		m.WhenReturn(v, m)
	}
	return m
}

func (m *QueryMock) Sort(fields ...string) IQuery {
	f := make([]interface{}, len(fields))
	for i, v := range fields {
		f[i] = v
	}
	return m.returnQuery("Sort", f...)
}

func (m *QueryMock) Skip(n int) IQuery {
	return m.returnQuery("Skip", n)
}

func (m *QueryMock) Limit(n int) IQuery {
	return m.returnQuery("Limit", n)
}

func (m *QueryMock) Page(page ...*pages.Page) IQuery {
	p := make([]interface{}, len(page))
	for i, v := range page {
		p[i] = v
	}
	return m.returnQuery("Page", p...)
}

func (m *QueryMock) WrapPage(result interface{}, page ...*pages.Page) (*pages.Paginated, error) {
	args := []interface{}{result}
	for _, v := range page {
		args = append(args, v)
	}
	ret, err := m.returnSingleWithError("WrapPage", args...)

	if ret != nil {
		return ret.(*pages.Paginated), err
	}

	return nil, err
}

func (m *QueryMock) Count() (int, error) {
	ret, err := m.returnSingleWithError("Count")

	if ret != nil {
		return ret.(int), err
	}

	return 0, err
}

func (m *QueryMock) All(result interface{}) error {
	return m.returnError("All", result)
}
