package mgo

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// NewSession creates an instance of ISession with the given *mgo.Session if passed as an arg.
// Note: The ISession instance returned will not work without a valid *mgo.Session.
func NewSession(s ...*mgo.Session) ISession {
	if len(s) > 0 {
		return &session{Session: s[0]}
	}
	return &session{}
}

// ISession is an interface which matches the contract for the `Session` struct in `gopkg.in/mgo.v2` package. The function documentation has been narrowed from the original
// in `gopkg.in/mgo.v2`. For additional documentation, please refer to the `mgo.collection` in the `gopkg.in/mgo.v2` package.
type ISession interface {
	ISessionExtensions

	// LiveServers returns a list of server addresses which are currently known to be alive.
	LiveServers() (addrs []string)

	// DB returns a value representing the named Database. If name is empty, the Database name provided in the dialed URL is used instead. If that is also empty, "test" is used as a
	// fallback in a way equivalent to the mongo shell.
	DB(name string) IDatabase

	// Login authenticates with MongoDB using the provided credential. The authentication is valid for the whole Session and will stay valid until Logout is explicitly called for
	// the same Database, or the Session is closed.
	Login(cred *mgo.Credential) error

	// LogoutAll removes all established authentication credentials for the Session.
	LogoutAll()

	// ResetIndexCache() clears the cache of previously ensured indexes. Following requests to EnsureIndex will contact the server.
	ResetIndexCache()

	// New creates a new Session with the same parameters as the original Session, including consistency, batch size, prefetching, safety mode, etc. The returned Session will use
	// sockets from the pool, so there's a chance that writes just performed in another Session may not yet be visible.
	//
	// Login information from the original Session will not be copied over into the new Session unless it was provided through the initial URL for the Dial function.
	New() ISession

	// Copy works just like New, but preserves the exact authentication information from the original Session.
	Copy() ISession

	// Clone works just like Copy, but also reuses the same socket as the original Session, in case it had already reserved one due to its consistency guarantees. This behavior
	// ensures that writes performed in the old Session are necessarily observed when using the new Session, as long as it was a strong or monotonic Session. That said, it also
	// means that long operations may cause other goroutines using the original Session to wait.
	Clone() ISession

	// Close terminates the Session.  It's a runtime error to use a Session after it has been closed.
	Close()

	// Refresh puts back any reserved sockets in use and restarts the consistency guarantees according to the current consistency setting for the Session.
	Refresh()

	// SetMode changes the consistency mode for the Session. The default mode is Strong.
	//     - See the SetMode documentation in `gopkg.in/mgo.v2` for more information.
	SetMode(consistency mgo.Mode, refresh bool)

	// Mode returns the current consistency mode for the Session.
	Mode() mgo.Mode

	// SetSyncTimeout sets the amount of time an operation with this Session will wait before returning an error in case a connection to a usable server can't be established. Set it
	// to zero to wait forever. The default value is 7 seconds.
	SetSyncTimeout(d time.Duration)

	// SetSocketTimeout sets the amount of time to wait for a non-responding socket to the Database before it is forcefully closed. The default timeout is 1 minute.
	SetSocketTimeout(d time.Duration)

	// SetCursorTimeout changes the standard timeout period that the server enforces on created cursors. The only supported value right now is 0, which disables the timeout.
	// The standard server timeout is 10 minutes.
	SetCursorTimeout(d time.Duration)

	// SetPoolLimit sets the maximum number of sockets in use in a single server before this Session will block waiting for a socket to be available. The default limit is 4096.
	//     - See the SetPoolLimit documentation in `gopkg.in/mgo.v2` for more information.
	SetPoolLimit(limit int)

	// SetBypassValidation sets whether the server should bypass the registered validation expressions executed when documents are inserted or modified, in the interest of preserving
	// invariants in the collection being modified. The default is to not bypass, and thus to perform the validation expressions registered for modified collections.
	SetBypassValidation(bypass bool)

	// SetBatch sets the default batch size used when fetching documents from the Database. It's possible to change this setting on a per-query basis as well, using the query.Batch method.
	//     - See the SetBatch documentation in `gopkg.in/mgo.v2` for more information.
	SetBatch(n int)

	// SetPrefetch sets the default point at which the next batch of results will be requested. When there are p*batch_size remaining documents cached in an Iter, the next batch
	// will be requested in background. The default prefetch value is 0.25.
	SetPrefetch(p float64)

	// Safe returns the current safety mode for the Session.
	Safe() (safe *mgo.Safe)

	// SetSafe changes the Session safety mode.
	//     - See the SetSafe documentation in `gopkg.in/mgo.v2` for more information.
	SetSafe(safe *mgo.Safe)

	// EnsureSafe compares the provided safety parameters with the ones currently in use by the Session and picks the most conservative choice for each setting.
	//     - See the EnsureSafe documentation in `gopkg.in/mgo.v2` for more information.
	EnsureSafe(safe *mgo.Safe)

	// Run issues the provided command on the "admin" Database and and unmarshals its result in the respective argument. The cmd argument may be either a string with the command name
	// itself, in which case an empty document of the form bson.M{cmd: 1} will be used, or it may be a full command document.
	//     - See the Session.Run documentation in `gopkg.in/mgo.v2` for more information.
	Run(cmd interface{}, result interface{}) error

	// SelectServers restricts communication to servers configured with the given tags. For example, the following statement restricts servers used for reading operations to those
	// with both tag "disk" set to "ssd" and tag "rack" set to 1:
	//     - See the Session.SelectServers documentation in `gopkg.in/mgo.v2` for more information.
	SelectServers(tags ...bson.D)

	// Ping runs a trivial ping command just to get in touch with the server.
	Ping() error

	// Fsync flushes in-memory writes to disk on the server the Session is established with. If async is true, the call returns immediately, otherwise it returns after the flush has
	// been made.
	Fsync(async bool) error

	// FsyncLock locks all writes in the specific server the Session is established with and returns. Any writes attempted to the server after it is successfully locked will block
	// until FsyncUnlock is called for the same server.
	//     - See the Session.FsyncLock documentation in `gopkg.in/mgo.v2` for more information.
	FsyncLock() error

	// FsyncUnlock releases the server for writes. See FsyncLock for details.
	FsyncUnlock() error

	// FindRef returns a query that looks for the document in the provided reference. For a DBRef to be resolved correctly at the Session level it must necessarily have the optional
	// DB field defined.
	FindRef(ref *mgo.DBRef) IQuery

	// DatabaseNames returns the names of non-empty databases present in the cluster.
	DatabaseNames() (names []string, err error)

	// BuildInfo retrieves the version and other details about the running MongoDB server.
	BuildInfo() (info mgo.BuildInfo, err error)

	// Returns the internal mgo.Session used by this implementation.
	S() *mgo.Session
}

// Session is the default implementation of ISession
type session struct {
	*mgo.Session
}

func (s *session) S() *mgo.Session {
	return s.Session
}

func (s *session) DB(name string) IDatabase {
	return fromDatabase(s.S().DB(name))
}

func (s *session) New() ISession {
	return fromSession(s.S().New())
}

func (s *session) Copy() ISession {
	return fromSession(s.S().Copy())
}

func (s *session) Clone() ISession {
	return fromSession(s.S().Clone())
}

func (s *session) FindRef(ref *mgo.DBRef) IQuery {
	return fromQuery(s.S().FindRef(ref))
}

func (s *session) LogoutAll() {
	println("SUCCESS!!!")
}

func fromSession(s *mgo.Session) ISession {
	if s == nil {
		return nil
	}
	return &session{Session: s}
}
