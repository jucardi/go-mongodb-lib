package mgo

import (
	"github.com/jucardi/go-mongodb-lib/log"
)

// ICollectionExtensions encapsulates the new extended functions to the original ICollection
type ICollectionExtensions interface {
	// MustEnsureIndex ensures an index with the given key exists, creating it with
	// the provided parameters if necessary.  If the index fails then this call
	// exists as a Fatal
	//
	//   {index} - the index in Bson form
	//
	MustEnsureIndex(index Index)

	// BulkUpsert allows multiple Upsert operations. Queues up the provided pairs of upserting instructions.
	// The first element of each pair selects which documents must be updated, and the second element defines how to update it.
	// Each pair matches exactly one document for updating at most.
	//
	// Enhanced to use bulk operations in the length of documents is more than the allowed 1000.
	BulkUpsert(pairs ...interface{}) (*BulkResult, error)
}

// MustEnsureIndex ensures an index with the given key exists, creating it with
// the provided parameters if necessary.  If the index fails then this call
// exists as a Fatal
//
//   {index} - the index in Bson form
//
func (c *collection) MustEnsureIndex(index Index) {
	if err := c.EnsureIndex(index); err != nil {
		log.Get().Fatal(err)
	} else {
		log.Get().Infof("collection [%s] index is up to date", c.Name)
	}
}

// Insert **Override of mgo.collection.Insert** inserts one or more documents in the respective collection.
// The override behavior converts the insert into a bulk operation if the length of documents is more than the allowed 1000 by MongoDB.
func (c *collection) Insert(docs ...interface{}) error {
	if len(docs) < mgoLim {
		return c.C().Insert(docs...)
	}
	_, err := NewBulk(c).Insert(docs...).Run()
	return err
}

// BulkUpsert allows multiple Upsert operations. Queues up the provided pairs of upserting instructions.
// The first element of each pair selects which documents must be updated, and the second element defines how to update it.
// Each pair matches exactly one document for updating at most.
//
// Enhanced to use bulk operations in the length of documents is more than the allowed 1000.
func (c *collection) BulkUpsert(pairs ...interface{}) (*BulkResult, error) {
	return NewBulk(c).Upsert(pairs...).Run()
}
