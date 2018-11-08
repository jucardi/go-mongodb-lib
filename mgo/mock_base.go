package mgo

import (
	"fmt"

	"github.com/jucardi/go-mongodb-lib/log"
	"github.com/jucardi/go-mongodb-lib/testutils"
)

type MockBase struct {
	times map[string]int
	when  map[string]testutils.WhenHandler
}

func newMock() *MockBase {
	return &MockBase{
		times: map[string]int{},
		when:  map[string]testutils.WhenHandler{},
	}
}

// When indicates what the expected behavior should be when a function is invoked.
func (m *MockBase) When(funcName string, f testutils.WhenHandler) {
	m.when[funcName] = func(args ...interface{}) []interface{} {
		log.Get().Debug(funcName, " func invoked")
		m.times[funcName]++
		return f(args...)
	}
}

// WhenReturn allows to set the return args without the need of a WhenHandler.
func (m *MockBase) WhenReturn(funcName string, retArgs ...interface{}) {
	m.When(funcName, func(args ...interface{}) []interface{} {
		return retArgs
	})
}

// Times asserts that the amount of times a function was invoked matches the provided 'expected'.
func (m *MockBase) Times(funcName string) int {
	return m.times[funcName]
}

// BulkTimes same as 'Times', but verifies multiple function calls in one call.
func (m *MockBase) BulkTimes(names []string, expected []int) error {
	if len(expected) != len(names) {
		panic("The arrays used in BulkTime must be the same size.")
	}

	for i, v := range names {
		if expected[i] != m.Times(v) {
			return fmt.Errorf("field %s times mismatch, expected %d, got %d", v, expected[i], m.Times(v))
		}
	}
	return nil
}

func (m *MockBase) exec(name string, args ...interface{}) []interface{} {
	if f, ok := m.when[name]; ok {
		return f(args...)
	}
	m.times[name]++
	return nil
}

func (m *MockBase) returnSingle(name string, args ...interface{}) interface{} {
	ret := m.exec(name, args...)
	if len(ret) > 0 && ret[0] != nil {
		return ret[0]
	}
	return nil
}

func (m *MockBase) returnDouble(name string, args ...interface{}) (interface{}, interface{}) {
	ret := m.exec(name, args...)

	if len(ret) < 2 {
		log.Get().Warn("Expected 2 returns for '%s', found %d", name, len(ret))
		return nil, nil
	}

	return ret[0], ret[1]
}

func (m *MockBase) returnError(name string, args ...interface{}) error {
	if val, ok := m.returnSingle(name, args...).(error); ok {
		return val
	}
	return nil
}

func (m *MockBase) returnSingleWithError(name string, args ...interface{}) (interface{}, error) {
	ret, err := m.returnDouble(name, args...)

	if err != nil {
		return ret, err.(error)
	}

	return ret, nil
}

func (m *MockBase) returnQuery(name string, args ...interface{}) IQuery {
	if val, ok := m.returnSingle(name, args...).(IQuery); ok {
		return val
	}
	return nil
}

func (m *MockBase) returnDB(name string, args ...interface{}) IDatabase {
	if val, ok := m.returnSingle(name, args...).(IDatabase); ok {
		return val
	}
	return nil
}

func (m *MockBase) returnSession(name string, args ...interface{}) ISession {
	if val, ok := m.returnSingle(name, args...).(ISession); ok {
		return val
	}
	return nil
}

func (m *MockBase) returnCollection(name string, args ...interface{}) ICollection {
	if val, ok := m.returnSingle(name, args...).(ICollection); ok {
		return val
	}
	return nil
}
