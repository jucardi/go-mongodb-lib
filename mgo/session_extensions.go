package mgo

import (
	"gopkg.in/mgo.v2"
)

// ISessionExtensions encapsulates the new extended functions to the original ISession
type ISessionExtensions interface {
	// SetDefaultSafe invokes SetSafe with the default safety instance.
	SetDefaultSafe()
}

func (s *session) SetDefaultSafe() {
	s.S().SetSafe(&mgo.Safe{})
}
