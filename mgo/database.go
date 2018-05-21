package mgo

import "gopkg.in/mgo.v2"

// NewDatabase creates an instance of IDatabase with the given *mgo.Database if passed as an arg.
// Note: The IDatabase instance returned will not work without a valid *mgo.Database.
func NewDatabase(db ...*mgo.Database) IDatabase {
	if len(db) > 0 {
		return &database{Database: db[0]}
	}
	return &database{}
}

// IDatabase is an interface which matches the contract for the `Database` struct in `gopkg.in/mgo.v2` package. The function documentation has been narrowed from the original
// in `gopkg.in/mgo.v2`. For additional documentation, please refer to the `mgo.collection` in the `gopkg.in/mgo.v2` package.
type IDatabase interface {
	IDatabaseExtensions

	// C returns a value representing the named collection.
	C(name string) ICollection

	// With returns a copy of db that uses Session s.
	With(s ISession) IDatabase

	// GridFS returns a GridFS value representing collections in db that follow the standard GridFS specification.
	GridFS(prefix string) *mgo.GridFS

	// Run issues the provided command on the db Database and unmarshals its result in the respective argument. The cmd argument may be either a string with the command name itself,
	// in which case an empty document of the form bson.M{cmd: 1} will be used, or it may be a full command document.
	Run(cmd interface{}, result interface{}) error

	// Login authenticates with MongoDB using the provided credential. The authentication is valid for the whole Session and will stay valid until Logout is explicitly called for
	// the same Database, or the Session is closed.
	Login(user, pass string) error

	// Logout removes any established authentication credentials for the Database.
	Logout()

	// UpsertUser updates the authentication credentials and the roles for a MongoDB user within the db Database. If the named user doesn't exist it will be created.
	// This method should only be used from MongoDB 2.4 and on. For older MongoDB releases, use the obsolete AddUser method instead.
	UpsertUser(user *mgo.User) error

	// AddUser creates or updates the authentication credentials of user within the db Database.
	// WARNING: This method is obsolete and should only be used with MongoDB 2.2 or earlier. For MongoDB 2.4 and on, use UpsertUser instead.
	AddUser(username, password string, readOnly bool) error

	// RemoveUser removes the authentication credentials of user from the Database.
	RemoveUser(user string) error

	// DropDatabase removes the entire Database including all of its collections.
	DropDatabase() error

	// FindRef returns a query that looks for the document in the provided reference. If the reference includes the DB field, the document will be retrieved from the respective Database.
	FindRef(ref *mgo.DBRef) IQuery

	// CollectionNames returns the collection names present in the db Database.
	CollectionNames() (names []string, err error)

	// Name returns the name of the database
	Name() string

	// Session returns the session used by the database
	Session() ISession

	// Returns the internal mgo.Database used by this implementation.
	DB() *mgo.Database
}

type database struct {
	*mgo.Database
}

func (d *database) DB() *mgo.Database {
	return d.Database
}

func (d *database) C(name string) ICollection {
	return fromCollection(d.DB().C(name))
}

func (d *database) With(s ISession) IDatabase {
	return fromDatabase(d.DB().With(s.S()))
}

func (d *database) FindRef(ref *mgo.DBRef) IQuery {
	return fromQuery(d.DB().FindRef(ref))
}

func (d *database) Name() string {
	return d.DB().Name
}

func (d *database) Session() ISession {
	return fromSession(d.DB().Session)
}

func fromDatabase(db *mgo.Database) IDatabase {
	if db == nil {
		return nil
	}
	return &database{Database: db}
}
