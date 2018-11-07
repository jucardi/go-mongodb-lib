package mgo

import (
	"gopkg.in/jucardi/go-streams.v1/streams"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// NewBulk creates an instance of ICollection with the given *mgo.Collection if passed as an arg.
// Note: The ICollection instance returned will not work without a valid *mgo.Collection.
func NewCollection(col ...*mgo.Collection) ICollection {
	if len(col) > 0 {
		return &collection{Collection: col[0]}
	}
	return &collection{}
}

// ICollection is an interface which matches the contract for the `collection` struct in `gopkg.in/mgo.v2` package. The function documentation has been narrowed from the original
// in `gopkg.in/mgo.v2`. For additional documentation, please refer to the `mgo.collection` in the `gopkg.in/mgo.v2` package.
type ICollection interface {
	// Set of extension functions that are not present in the original `mgo` package are defined in the following interface(s):
	ICollectionExtensions

	// With returns a copy of c that uses Session s.
	With(s ISession) ICollection

	// EnsureIndexKey ensures an index with the given key exists, creating it
	EnsureIndexKey(key ...string) error

	// EnsureIndex ensures an index with the given key exists, creating it with
	EnsureIndex(index Index) error

	// DropIndex drops the index with the provided key from the c collection.
	DropIndex(key ...string) error

	// DropIndexName removes the index with the provided index name.
	DropIndexName(name string) error

	// Indexes returns a list of all indexes for the collection.
	Indexes() (indexes []Index, err error)

	// Find prepares a query using the provided document. The document may be a map or a struct value capable of being marshalled with bson. The map may be a generic one using
	// interface{} for its key and/or values, such as bson.M, or it may be a properly typed map.  roviding nil as the document is equivalent to providing an empty document such
	// as bson.M{}.
	Find(query interface{}) IQuery

	// Repair returns an iterator that goes over all recovered documents in the collection, in a best-effort manner. This is most useful when there are damaged data files. Multiple
	// copies of the same document may be returned by the iterator.
	Repair() IIter

	// FindId is a convenience helper equivalent to:
	FindId(id interface{}) IQuery

	// pipe prepares a pipeline to aggregate. The pipeline document must be a slice built in terms of the aggregation framework language.
	Pipe(pipeline interface{}) IPipe

	// NewIter returns a newly created iterator with the provided parameters. Using this method is not recommended unless the desired functionality is not yet exposed via a more
	// convenient interface (Find, pipe, etc).
	NewIter(session ISession, firstBatch []bson.Raw, cursorId int64, err error) IIter

	// Insert inserts one or more documents in the respective collection.  In case the Session is in safe mode (see the SetSafe method) and an error happens while inserting the
	// provided documents, the returned error will be of type *LastError.
	Insert(docs ...interface{}) error

	// Update finds a single document matching the provided selector document and modifies it according to the update document.
	Update(selector interface{}, update interface{}) error

	// UpdateId is a convenience helper equivalent to:
	UpdateId(id interface{}, update interface{}) error

	// UpdateAll finds all documents matching the provided selector document and modifies them according to the update document.
	UpdateAll(selector interface{}, update interface{}) (info *ChangeInfo, err error)

	// Upsert finds a single document matching the provided selector document and modifies it according to the update document.  If no document matching the selector is found, the
	// update document is applied to the selector document and the result is inserted in the collection.
	Upsert(selector interface{}, update interface{}) (info *ChangeInfo, err error)

	// UpsertId is a convenience helper equivalent to:
	UpsertId(id interface{}, update interface{}) (info *ChangeInfo, err error)

	// Remove finds a single document matching the provided selector document and removes it from the Database.
	Remove(selector interface{}) error

	// RemoveId is a convenience helper equivalent to:
	RemoveId(id interface{}) error

	// RemoveAll finds all documents matching the provided selector document and removes them from the Database.
	RemoveAll(selector interface{}) (info *ChangeInfo, err error)

	// DropCollection removes the entire collection including all of its documents.
	DropCollection() error

	// Create explicitly creates the c collection with details of info. MongoDB creates collections automatically on use, so this method is only necessary when creating collection
	// with non-default characteristics, such as capped collections.
	Create(info *CollectionInfo) error

	// Count returns the total number of documents in the collection.
	Count() (n int, err error)

	// Database returns the database the collection belongs to
	Database() IDatabase

	// Name returns the name of the collection ("collection")
	Name() string

	// FullName returns the full name of the collection ("db.collection")
	FullName() string

	// Bulk returns a value to prepare the execution of a bulk operation.
	Bulk() IBulk

	// Returns the internal mgo.collection used by this implementation.
	C() *mgo.Collection
}

// collection is the default implementation of ICollection
type collection struct {
	*mgo.Collection
}

func (c *collection) C() *mgo.Collection {
	return c.Collection
}
func (c *collection) With(s ISession) ICollection {
	return c.update(c.C().With(s.S()))
}

func (c *collection) Find(query interface{}) IQuery {
	return fromQuery(c.C().Find(query))
}

func (c *collection) FindId(id interface{}) IQuery {
	return fromQuery(c.C().FindId(id))
}

func (c *collection) Create(info *CollectionInfo) error {
	i := mgo.CollectionInfo(*info)
	return c.C().Create(&i)
}

func (c *collection) RemoveAll(selector interface{}) (*ChangeInfo, error) {
	info, err := c.C().RemoveAll(selector)
	return makeChangeInfo(info), err
}

func (c *collection) UpsertId(id interface{}, update interface{}) (*ChangeInfo, error) {
	info, err := c.C().UpsertId(id, update)
	return makeChangeInfo(info), err
}

func (c *collection) Upsert(selector interface{}, update interface{}) (*ChangeInfo, error) {
	info, err := c.C().Upsert(selector, update)
	return makeChangeInfo(info), err
}

func (c *collection) UpdateAll(selector interface{}, update interface{}) (*ChangeInfo, error) {
	info, err := c.C().UpdateAll(selector, update)
	return makeChangeInfo(info), err
}

func (c *collection) NewIter(session ISession, firstBatch []bson.Raw, cursorId int64, err error) IIter {
	return c.C().NewIter(session.S(), firstBatch, cursorId, err)
}

func (c *collection) Pipe(pipeline interface{}) IPipe {
	return fromPipe(c.C().Pipe(pipeline))
}

func (c *collection) Repair() IIter {
	return c.C().Repair()
}

func (c *collection) Indexes() ([]Index, error) {
	ret, err := c.C().Indexes()
	return streams.From(ret).Map(func(x interface{}) interface{} {
		return Index(x.(mgo.Index))
	}).ToArray().([]Index), err
}

func (c *collection) EnsureIndex(index Index) error {
	return c.C().EnsureIndex(mgo.Index(index))
}

func (c *collection) Bulk() IBulk {
	return NewBulk(c)
}

func (c *collection) Database() IDatabase {
	return fromDatabase(c.C().Database)
}

func (c *collection) Name() string {
	return c.C().Name
}

func (c *collection) FullName() string {
	return c.C().FullName
}

func (c *collection) update(col *mgo.Collection) ICollection {
	c.Collection = col
	return c
}

func fromCollection(col *mgo.Collection) ICollection {
	if col == nil {
		return nil
	}
	return &collection{Collection: col}
}
