package mgo

import (
	"time"

	"github.com/jucardi/go-mongodb-lib/log"
	"gopkg.in/mgo.v2"
)

var (
	// ErrNotFound is the error returned when no results are found in a mongo operation.
	ErrNotFound = mgo.ErrNotFound

	// ErrCursor is the error returned when the cursor used in a mongo operation is not valid.
	ErrCursor = mgo.ErrCursor
)

// Dial establishes a new session to the cluster identified by the given seed
// server(s). The session will enable communication with all of the servers in
// the cluster, so the seed servers are used only to find out about the cluster
// topology.
//
//     - See mgo.Dial documentation in `gopkg.in/mgo.v2` for more information.
//
func Dial(url string) (ISession, error) {
	s, err := mgo.Dial(url)

	for i := 1; err != nil && i <= Config().DialMaxRetries; i++ {
		log.Get().Errorf("Can't connect to mongo on '%s': %v. Retrying in %v", url, err, Config().DialRetryTimeout)
		time.Sleep(Config().DialRetryTimeout)
		log.Get().Warnf("Retrying to connect to mongo, attempt %d of %d", i, Config().DialMaxRetries)
		s, err = mgo.Dial(url)
	}

	return fromSession(s), err
}

// DialWithTimeout works like Dial, but uses timeout as the amount of time to
// wait for a server to respond when first connecting and also on follow up
// operations in the session. If timeout is zero, the call may block
// forever waiting for a connection to be made.
//
// See SetSyncTimeout for customizing the timeout for the session.
func DialWithTimeout(url string, timeout time.Duration) (ISession, error) {
	s, err := mgo.DialWithTimeout(url, timeout)
	return fromSession(s), err
}

// DialWithInfo establishes a new session to the cluster identified by info.
func DialWithInfo(info *mgo.DialInfo) (ISession, error) {
	s, err := mgo.DialWithInfo(info)
	return fromSession(s), err
}

// IsDup returns whether err informs of a duplicate key error because
// a primary key index or a secondary unique index already has an entry
// with the given value.
func IsDup(err error) bool {
	return mgo.IsDup(err)
}
