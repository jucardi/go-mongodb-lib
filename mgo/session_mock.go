package mgo

import "gopkg.in/mgo.v2"

var sessionFuncs = []string{"New", "Copy", "Clone"}

// QueryMock is a mock implementation of IQuery
type SessionMock struct {
	ISession
	*MockBase
}

// MockSession returns a new instance of IQuery for mocking purposes
//
//   {q}:   An optional instance of IQuery to handle real calls if wanted.
//
func MockSession(q ...ISession) *SessionMock {
	ret := &SessionMock{
		ISession: NewSession(),
		MockBase: newMock(),
	}
	if len(q) > 0 {
		ret.ISession = q[0]
	}
	return ret
}

func (s *SessionMock) init() *SessionMock {
	// Sets the default return value for the methods that return IQuery. By default they'll return the mock.
	for _, v := range sessionFuncs {
		s.WhenReturn(v, s)
	}
	return s
}

func (s *SessionMock) New() ISession {
	return s.returnSession("New")
}

func (s *SessionMock) Copy() ISession {
	return s.returnSession("Copy")
}

func (s *SessionMock) Clone() ISession {
	return s.returnSession("Clone")
}

func (s *SessionMock) DB(name string) IDatabase {
	return s.returnDB("DB", name)
}

func (s *SessionMock) FindRef(ref *mgo.DBRef) IQuery {
	return s.returnQuery("FindRef", ref)
}
