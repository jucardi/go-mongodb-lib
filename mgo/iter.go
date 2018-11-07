package mgo

type IIter interface {
	// Err returns nil if no errors happened during iteration, or the actual
	// error otherwise.
	//     - See the Err documentation in `gopkg.in/mgo.v2` for more information.
	Err() error

	// Close kills the server cursor used by the iterator, if any, and returns
	// nil if no errors happened during iteration, or the actual error otherwise.
	//     - See the Close documentation in `gopkg.in/mgo.v2` for more information.
	Close() error

	// Done returns true only if a follow up Next call is guaranteed
	// to return false.
	//     - See the Done documentation in `gopkg.in/mgo.v2` for more information.
	Done() bool

	// Timeout returns true if Next returned false due to a timeout of
	// a tailable cursor. In those cases, Next may be called again to continue
	// the iteration at the previous cursor position.
	Timeout() bool

	// Next retrieves the next document from the result set, blocking if necessary.
	// This method will also automatically retrieve another batch of documents from
	// the server when the current one is exhausted, or before that in background
	// if pre-fetching is enabled (see the query.Prefetch and Session.SetPrefetch
	// methods).
	//     - See the Next documentation in `gopkg.in/mgo.v2` for more information.
	Next(result interface{}) bool

	// All retrieves all documents from the result set into the provided slice
	// and closes the iterator.
	//     - See the All documentation in `gopkg.in/mgo.v2` for more information.
	All(result interface{}) error

	// The For method is obsolete and will be removed in a future release.
	// See Iter as an elegant replacement.
	For(result interface{}, f func() error) (err error)
}
