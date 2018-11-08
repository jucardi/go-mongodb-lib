package mgo

// TODO: Improve functionality like QueryMock

// CollectionMock: Is the mock struct use for ICollection mocking.
type CollectionMock struct {
	ICollection
	findMocks   map[string]IQuery
	insertMocks map[string]error
	times       map[string]int
}

// MockCollection returns a new instance of ICollection for mocking purposes
//
//   {t}:    The instance of *testing.T used in the test
//   {col}:  An optional instance of ICollection to handle real calls if wanted.
//
func MockCollection(col ...ICollection) *CollectionMock {
	ret := &CollectionMock{
		findMocks:   make(map[string]IQuery),
		insertMocks: make(map[string]error),
		times:       make(map[string]int),
	}
	if len(col) > 0 {
		ret.ICollection = col[0]
	}
	return ret
}

func (m *CollectionMock) WhenFind(query interface{}, output IQuery) {
	m.findMocks[getKey(query)] = output
}

func (m *CollectionMock) WhenInsert(docs interface{}, output error) {
	m.insertMocks[getKey(docs)] = output
}

func (m *CollectionMock) Times(funcName string) int {
	return m.times[funcName]
}

func (m *CollectionMock) Clear(funcName string) {
	for k := range m.times {
		delete(m.times, k)
	}
	for k := range m.insertMocks {
		delete(m.insertMocks, k)
	}
	for k := range m.findMocks {
		delete(m.findMocks, k)
	}
}

func (m *CollectionMock) Find(query interface{}) IQuery {
	m.times["Find"]++
	if _, ok := m.findMocks[getKey(query)]; ok {
		return m.findMocks[getKey(query)]
	} else if _, ok := m.findMocks["any"]; ok {
		return m.findMocks["any"]
	}
	return nil
}

func (m *CollectionMock) Insert(docs ...interface{}) error {
	m.times["Insert"]++
	if _, ok := m.insertMocks[getKey(docs)]; ok {
		return m.insertMocks[getKey(docs)]
	} else if _, ok := m.insertMocks["any"]; ok {
		return m.insertMocks["any"]
	}

	return nil
}
