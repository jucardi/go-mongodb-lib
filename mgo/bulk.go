package mgo

import (
	"math"
	"reflect"

	"gopkg.in/mgo.v2"
)

// Indicates the max amount of items a mongo operation can handle.
const mgoLim = 1000

// NewBulk creates an instance of IBulk with the given ICollection if passed as an arg.
// Note: The IBulk instance returned will not work without a valid ICollection.
func NewBulk(col ...ICollection) IBulk {
	if len(col) > 0 {
		return &bulk{
			col:     col[0].C(),
			ordered: true,
		}
	}
	return &bulk{
		ordered: true,
	}
}

// IBulk is an interface which matches the contract for the `Bulk` struct in `gopkg.in/mgo.v2` package. The function documentation has been narrowed from the original
// in `gopkg.in/mgo.v2`. For additional documentation, please refer to the `mgo.Bulk` in the `gopkg.in/mgo.v2` package.
type IBulk interface {
	// Unordered puts the bulk operation in unordered mode.
	//
	// In unordered mode the individual operations may be sent
	// out of order, which means latter operations may proceed
	// even if prior ones have failed.
	Unordered() IBulk

	// Insert queues up the provided documents for insertion.
	Insert(docs ...interface{}) IBulk

	// Remove queues up the provided selectors for removing matching documents.
	// Each selector will remove only a single matching document.
	Remove(selectors ...interface{}) IBulk

	// RemoveAll queues up the provided selectors for removing all matching documents.
	// Each selector will remove all matching documents.
	RemoveAll(selectors ...interface{}) IBulk

	// Update queues up the provided pairs of updating instructions.
	// The first element of each pair selects which documents must be
	// updated, and the second element defines how to update it.
	// Each pair matches exactly one document for updating at most.
	Update(pairs ...interface{}) IBulk

	// UpdateAll queues up the provided pairs of updating instructions.
	// The first element of each pair selects which documents must be
	// updated, and the second element defines how to update it.
	// Each pair updates all documents matching the selector.
	UpdateAll(pairs ...interface{}) IBulk

	// Upsert queues up the provided pairs of upserting instructions.
	// The first element of each pair selects which documents must be
	// updated, and the second element defines how to update it.
	// Each pair matches exactly one document for updating at most.
	Upsert(pairs ...interface{}) IBulk

	// Run runs all the operations queued up.
	//
	// If an error is reported on an unordered bulk operation, the error value may
	// be an aggregation of all issues observed. As an exception to that, Insert
	// operations running on MongoDB versions prior to 2.6 will report the last
	// error only due to a limitation in the wire protocol.
	Run() (*BulkResult, error)
}

// Bulk is the default implementation of IBulk
type bulk struct {
	bulks   []*mgo.Bulk
	col     *mgo.Collection
	ordered bool
}

func (b *bulk) Unordered() IBulk {
	b.ordered = false
	return b
}

func (b *bulk) Insert(docs ...interface{}) IBulk {
	return b.add("Insert", docs)
}

func (b *bulk) Remove(selectors ...interface{}) IBulk {
	return b.add("Remove", selectors)
}

func (b *bulk) RemoveAll(selectors ...interface{}) IBulk {
	return b.add("RemoveAll", selectors)
}

func (b *bulk) Update(pairs ...interface{}) IBulk {
	return b.add("Update", pairs)
}

func (b *bulk) UpdateAll(pairs ...interface{}) IBulk {
	return b.add("UpdateAll", pairs)
}

func (b *bulk) Upsert(pairs ...interface{}) IBulk {
	return b.add("Upsert", pairs)
}

func (b *bulk) Run() (*BulkResult, error) {
	ret := &BulkResult{}

	for _, bulk := range b.bulks {
		if !b.ordered {
			bulk.Unordered()
		}

		r, err := bulk.Run()

		if r != nil {
			ret.Modified += r.Modified
			ret.Matched += r.Matched
		}

		if err != nil {
			return ret, err
		}
	}

	return ret, nil
}

func (b *bulk) add(f string, items []interface{}) IBulk {
	l := float64(len(items))
	lim := float64(mgoLim)

	for i := 0; i+1 <= int(math.Ceil(l/lim)); i++ {
		blk := b.col.Bulk()
		top := int(math.Max(l, float64((i+1)*mgoLim)))

		val := reflect.ValueOf(blk)
		fn := val.MethodByName(f)
		arg := items[i*mgoLim : top]
		fn.Call([]reflect.Value{reflect.ValueOf(arg)})
		b.bulks = append(b.bulks, blk)
	}
	return b
}
