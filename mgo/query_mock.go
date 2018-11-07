package mgo

import (
	"fmt"
	"testing"

	"github.com/jucardi/go-mongodb-lib/pages"
	"github.com/jucardi/go-mongodb-lib/testutils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// TODO: Handle conditionals on arguments to support more use cases when testing. Add all IQuery methods.

// All functions that return IQuery to easily initialize
var queryFuncs = []string{"Batch", "Prefetch", "Skip", "Limit", "Select", "Sort", "Count", "Explain", "Hint", "SetMaxScan", "SetMaxTime", "Snapshot", "Comment", "LogReplay", "Page"}

// QueryMock is a mock implementation of IQuery
type QueryMock struct {
	IQuery
	t     *testing.T
	times map[string]int
	when  map[string]WhenHandler
}

// MockQuery returns a new instance of IQuery for mocking purposes
//
//   {t}:   The instance of *testing.T used in the test
//   {q}:   An optional instance of IQuery to handle real calls if wanted.
//
func MockQuery(t *testing.T, q ...IQuery) *QueryMock {
	ret := &QueryMock{
		IQuery: NewQuery(),
		t:      t,
	}
	if len(q) > 0 {
		ret.IQuery = q[0]
	}
	return ret.init()
}

func (m *QueryMock) init() *QueryMock {
	m.times = make(map[string]int)
	m.when = make(map[string]WhenHandler)

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
	return m.doToQuery("Sort", f...)
}

func (m *QueryMock) Skip(n int) IQuery {
	return m.doToQuery("Skip", n)
}

func (m *QueryMock) Limit(n int) IQuery {
	return m.doToQuery("Limit", n)
}

func (m *QueryMock) Page(page ...*pages.Page) IQuery {
	p := make([]interface{}, len(page))
	for i, v := range page {
		p[i] = v
	}
	return m.doToQuery("Page", p...)
}

func (m *QueryMock) WrapPage(result interface{}, page ...*pages.Page) (*pages.Paginated, error) {
	args := []interface{}{result}
	for _, v := range page {
		args = append(args, v)
	}
	ret, err := m.doTwoReturnsError("WrapPage", args...)

	if ret != nil {
		return ret.(*pages.Paginated), err
	}

	return nil, err
}

func (m *QueryMock) Count() (int, error) {
	ret, err := m.doTwoReturnsError("Count")

	if ret != nil {
		return ret.(int), err
	}

	return 0, err
}

func (m *QueryMock) All(result interface{}) error {
	return m.doToError("All", result)
}

// Times: Asserts that the amount of times a function was invoked matches the provided 'expected'.
func (m *QueryMock) Times(funcName string, expected int) {
	assert.Equal(m.t, expected, m.times[funcName], "Mismatch count for method '%s'", funcName)
}

// BulkTimes: Same as 'Times', but verifies multiple function calls in one call.
func (m *QueryMock) BulkTimes(names []string, expected []int) {
	if len(expected) != len(names) {
		panic("The arrays used in BulkTime must be the same size.")
	}

	for i, v := range names {
		m.Times(v, expected[i])
	}
}

// When: Indicates what the expected behavior should be when a function is invoked.
func (m *QueryMock) When(funcName string, f WhenHandler) {
	m.when[funcName] = func(t *testing.T, args ...interface{}) []interface{} {
		logrus.Debug(funcName, " func invoked")
		m.times[funcName]++
		return f(t, args...)
	}
}

// WhenReturn: Allows to set the return args without the need of a WhenHandler.
func (m *QueryMock) WhenReturn(funcName string, retArgs ...interface{}) {
	m.When(funcName, func(t *testing.T, args ...interface{}) []interface{} {
		return testutils.MakeReturn(retArgs...)
	})
}

func (m *QueryMock) do(name string, args ...interface{}) []interface{} {
	if f, ok := m.when[name]; ok {
		return f(m.t, args...)
	}
	m.times[name]++
	return nil
}

func (m *QueryMock) doToQuery(name string, args ...interface{}) IQuery {
	ret := m.do(name, args...)
	if len(ret) > 0 && ret[0] != nil {
		return ret[0].(IQuery)
	}
	return nil
}

func (m *QueryMock) doToError(name string, args ...interface{}) error {
	ret := m.do(name, args...)
	if len(ret) > 0 && ret[0] != nil {
		return ret[0].(error)
	}
	return nil
}

func (m *QueryMock) doTwoReturnsError(name string, args ...interface{}) (interface{}, error) {
	ret := m.do(name, args...)

	if len(ret) < 2 {
		fmt.Printf("\nWARN! Expected 2 returns for '%s', found %d\n\n", name, len(ret))
		return nil, nil
	}

	if ret[1] != nil {
		return ret[0], ret[1].(error)
	}
	return ret[0], nil
}
