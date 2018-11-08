package mgo

import (
	"encoding/json"
)

// TODO: Improve functionality like QueryMock

// DatabaseMock is a mock of IDatabase
type DatabaseMock struct {
	IDatabase
	collections map[string]ICollection
	runMocks    map[string]error
	times       map[string]int
}

// MockDb returns a new instance of IDatabase for mocking purposes
//
//   {db}:   An optional instance of IDatabase to handle real calls if wanted.
//
func MockDb(db ...IDatabase) *DatabaseMock {
	ret := &DatabaseMock{
		collections: make(map[string]ICollection),
		runMocks:    make(map[string]error),
		times:       make(map[string]int),
	}
	if len(db) > 0 {
		ret.IDatabase = db[0]
	}
	return ret
}

func (m *DatabaseMock) WhenC(colName string, output ICollection) {
	m.collections[colName] = output
}

func (m *DatabaseMock) WhenRun(cmd interface{}, output error) {
	m.runMocks[getKey(cmd)] = output
}

func (m *DatabaseMock) Times(funcName string) int {
	return m.times[funcName]
}

func (m *DatabaseMock) Clear(funcName string) {
	for k := range m.times {
		delete(m.times, k)
	}
	for k := range m.collections {
		delete(m.collections, k)
	}
	for k := range m.runMocks {
		delete(m.runMocks, k)
	}
}

func (m *DatabaseMock) C(collectionName string) ICollection {
	m.times["C"]++
	if _, ok := m.collections[collectionName]; ok {
		return m.collections[collectionName]
	}
	return nil
}

func (m *DatabaseMock) Run(cmd interface{}, result interface{}) error {
	m.times["Run"]++
	if _, ok := m.runMocks[getKey(cmd)]; ok {
		return m.runMocks[getKey(cmd)]
	} else if _, ok := m.runMocks["any"]; ok {
		return m.runMocks["any"]
	}
	return nil
}

func getKey(obj interface{}) string {
	if obj == nil {
		return "any"
	}

	b, _ := json.Marshal(obj)
	return string(b)
}
