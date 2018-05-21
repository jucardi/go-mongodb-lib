package mgo

import (
	"gopkg.in/mgo.v2"
	"time"
)

// NewQuery creates an instance of IQuery with the given *mgo.Query if passed as an arg.
// Note: The IQuery instance returned will not work without a valid *mgo.Query.
func NewQuery(q ...*mgo.Query) IQuery {
	if len(q) > 0 {
		return &query{Query: q[0]}
	}
	return &query{}
}

// IQuery is an interface which matches the contract for the `query` struct in `gopkg.in/mgo.v2` package. The function documentation has been narrowed from the original
// in `gopkg.in/mgo.v2`. For additional documentation, please refer to the `mgo.Collection` in the `gopkg.in/mgo.v2` package.
type IQuery interface {
	// Set of extension functions that are not present in the original `mgo` package are defined in the following interface(s):
	IQueryPageExtension

	// The default batch size is defined by the Database itself.  As of this writing, MongoDB will use an initial size of min(100 docs, 4MB) on the first batch, and 4MB on remaining ones.
	Batch(n int) IQuery

	// Prefetch sets the point at which the next batch of results will be requested. When there are p*batch_size remaining documents cached in an Iter, the next batch will be
	// requested in background. For instance, when using this:
	Prefetch(p float64) IQuery

	// Skip skips over the n initial documents from the query results. Note that this only makes sense with capped collections where documents are naturally ordered by insertion
	// time, or with sorted results.
	Skip(n int) IQuery

	// Limit restricts the maximum number of documents retrieved to n, and also changes the batch size to the same value.  Once n documents have been returned by Next, the following
	// call will return ErrNotFound.
	Limit(n int) IQuery

	// Select enables selecting which fields should be retrieved for the results found. For example, the following query would only retrieve the name field:
	Select(selector interface{}) IQuery

	// Sort asks the Database to order returned documents according to the provided field names. A field name may be prefixed by - (minus) for it to be sorted in reverse order.
	Sort(fields ...string) IQuery

	// Explain returns a number of details about how the MongoDB server would execute the requested query, such as the number of objects examined, the number of times the read lock
	// was yielded to allow writes to go in, and so on.
	Explain(result interface{}) error

	// Hint will include an explicit "hint" in the query to force the server to use a specified index, potentially improving performance in some situations. The provided parameters
	// are the fields that compose the key of the index to be used. For details on how the indexKey may be built, see the EnsureIndex method.
	Hint(indexKey ...string) IQuery

	// SetMaxScan constrains the query to stop after scanning the specified number of documents.
	SetMaxScan(n int) IQuery

	// SetMaxTime constrains the query to stop after running for the specified time. When the time limit is reached MongoDB automatically cancels the query.
	SetMaxTime(d time.Duration) IQuery

	// Snapshot will force the performed query to make use of an available index on the _id field to prevent the same document from being returned more than once in a single
	// iteration. This might happen without this setting in situations when the document changes in size and thus has to be moved while the iteration is running.
	Snapshot() IQuery

	// Comment adds a comment to the query to identify it in the Database profiler output.
	Comment(comment string) IQuery

	// LogReplay enables an option that optimizes queries that are typically made on the MongoDB oplog for replaying it. This is an internal implementation aspect and most likely
	// uninteresting for other uses. It has seen at least one use case, though, so it's exposed via the API.
	LogReplay() IQuery

	// One executes the query and unmarshals the first obtained document into the result argument. The result must be a struct or map value capable of being unmarshalled into by
	// gobson. This function blocks until either a result is available or an error happens.  For example:
	One(result interface{}) error

	// Iter executes the query and returns an iterator capable of going over all the results. Results will be returned in batches of configurable size (see the Batch method) and more
	// documents will be requested when a configurable number of documents is iterated over (see the Prefetch method).
	Iter() IIter

	// Tail returns a tailable iterator. Unlike a normal iterator, a tailable iterator may wait for new values to be inserted in the Collection once the end of the current result set
	// is reached, A tailable iterator may only be used with capped collections.
	//     - See the Tail documentation in `gopkg.in/mgo.v2` for more information.
	Tail(timeout time.Duration) IIter

	// All works like Iter.All.
	All(result interface{}) error

	// Count returns the total number of documents in the result set.
	Count() (n int, err error)

	// Distinct unmarshals into result the list of distinct values for the given key.
	Distinct(key string, result interface{}) error

	// MapReduce executes a map/reduce job for documents covered by the query. That kind of job is suitable for very flexible bulk aggregation of data performed at the server side
	// via Javascript functions.
	//     - See the MapReduce documentation in `gopkg.in/mgo.v2` for more information.
	MapReduce(job *MapReduce, result interface{}) (info *MapReduceInfo, err error)

	// Apply runs the findAndModify MongoDB command, which allows updating, upserting or removing a document matching a query and atomically returning either the old version (the
	// default) or the new version of the document (when ReturnNew is true). If no objects are found Apply returns ErrNotFound.
	//     - See the Apply documentation in `gopkg.in/mgo.v2` for more information.
	Apply(change Change, result interface{}) (info *ChangeInfo, err error)

	// Returns the internal mgo.query used by this implementation.
	Q() *mgo.Query
}

// query is the default implementation of IQuery
type query struct {
	*mgo.Query
}

func (q *query) Q() *mgo.Query {
	return q.Query
}

func (q *query) Batch(n int) IQuery {
	return q.update(q.Q().Batch(n))
}

func (q *query) Prefetch(p float64) IQuery {
	return q.update(q.Q().Prefetch(p))
}

func (q *query) Skip(n int) IQuery {
	return q.update(q.Q().Skip(n))
}

func (q *query) Limit(n int) IQuery {
	if n <= 0 {
		return q
	}
	return q.update(q.Q().Limit(n))
}

func (q *query) Select(selector interface{}) IQuery {
	return q.update(q.Q().Select(selector))
}

func (q *query) Sort(fields ...string) IQuery {
	if len(fields) == 0 {
	}
	return q.update(q.Q().Sort(fields...))
}

func (q *query) Hint(indexKey ...string) IQuery {
	return q.update(q.Q().Hint(indexKey...))
}

func (q *query) SetMaxScan(n int) IQuery {
	return q.update(q.Q().SetMaxScan(n))
}

func (q *query) SetMaxTime(d time.Duration) IQuery {
	return q.update(q.Q().SetMaxTime(d))
}

func (q *query) Snapshot() IQuery {
	return q.update(q.Q().Snapshot())
}

func (q *query) Comment(comment string) IQuery {
	return q.update(q.Q().Comment(comment))
}

func (q *query) LogReplay() IQuery {
	return q.update(q.Q().LogReplay())
}

func (q *query) Iter() IIter {
	return q.Q().Iter()
}

func (q *query) Tail(timeout time.Duration) IIter {
	return q.Q().Tail(timeout)
}

func (q *query) MapReduce(job *MapReduce, result interface{}) (*MapReduceInfo, error) {
	info, err := q.Q().MapReduce(makeMapReduce(job), result)
	return makeMapReduceInfo(info), err
}

func (q *query) Apply(change Change, result interface{}) (*ChangeInfo, error) {
	info, err := q.Q().Apply(mgo.Change(change), result)
	return makeChangeInfo(info), err
}

// update: Updates the inner *mgo.query contained by this instance.
func (q *query) update(query *mgo.Query) IQuery {
	q.Query = query
	return q
}

func fromQuery(q *mgo.Query) IQuery {
	return &query{Query: q}
}
