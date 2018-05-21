package mgo

// IDatabaseExtensions encapsulates the new extended functions to the original IDatabase
type IDatabaseExtensions interface {
	// MustEnsureIndex ensures an index with the given key exists, creating it with
	// the provided parameters if necessary.  If the index fails then this call
	// exists as a Fatal
	//
	// {index}      - the index in Bson form
	// {collection} - the name of the collection to apply the index to
	MustEnsureIndex(index Index, collection string)
}

// MustEnsureIndex ensures an index with the given key exists, creating it with
// the provided parameters if necessary.  If the index fails then this call
// exists as a Fatal
//
// {index}      - the index in Bson form
// {collection} - the name of the collection to apply the index to
func (d *database) MustEnsureIndex(index Index, collection string) {
	d.C(collection).MustEnsureIndex(index)
}
