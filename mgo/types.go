package mgo

import (
	"time"

	"gopkg.in/mgo.v2"
)

// To facilitate only importing one package when using these

// Change See the Change documentation in `gopkg.in/mgo.v2` for more information.
type Change struct {
	Update    interface{} // The update document
	Upsert    bool        // Whether to insert in case the document isn't found
	Remove    bool        // Whether to remove the document found rather than updating
	ReturnNew bool        // Should the modified document be returned rather than the old one
}

// ChangeInfo See the ChangeInfo documentation in `gopkg.in/mgo.v2` for more information.
type ChangeInfo struct {
	Updated    int
	Removed    int         // Number of documents removed
	Matched    int         // Number of documents matched but not necessarily changed
	UpsertedId interface{} // Upserted _id field, when not explicitly provided
}

// MapReduce See the MapReduce documentation in `gopkg.in/mgo.v2` for more information.
type MapReduce struct {
	Map      string      // Map Javascript function code (required)
	Reduce   string      // Reduce Javascript function code (required)
	Finalize string      // Finalize Javascript function code (optional)
	Out      interface{} // Output collection name or document. If nil, results are inlined into the result parameter.
	Scope    interface{} // Optional global scope for Javascript functions
	Verbose  bool
}

// MapReduceInfo See the MapReduceInfo documentation in `gopkg.in/mgo.v2` for more information.
type MapReduceInfo struct {
	InputCount  int                // Number of documents mapped
	EmitCount   int                // Number of times reduce called emit
	OutputCount int                // Number of documents in resulting collection
	Database    string             // Output database, if results are not inlined
	Collection  string             // Output collection, if results are not inlined
	Time        int64              // Time to run the job, in nanoseconds
	VerboseTime *mgo.MapReduceTime // Only defined if Verbose was true
}

// MapReduceTime See the MapReduceTime documentation in `gopkg.in/mgo.v2` for more information.
type MapReduceTime struct {
	Total    int64 // Total time, in nanoseconds
	Map      int64 // Time within map function, in nanoseconds
	EmitLoop int64 // Time within the emit/map loop, in nanoseconds
}

// CollectionInfo See the CollectionInfo documentation in `gopkg.in/mgo.v2` for more information.
type CollectionInfo struct {
	DisableIdIndex   bool
	ForceIdIndex     bool
	Capped           bool
	MaxBytes         int
	MaxDocs          int
	Validator        interface{}
	ValidationLevel  string
	ValidationAction string
	StorageEngine    interface{}
}

// Index See the Index documentation in `gopkg.in/mgo.v2` for more information.
type Index struct {
	Key              []string // Index key fields; prefix name with dash (-) for descending order
	Unique           bool     // Prevent two documents from having the same index key
	DropDups         bool     // Drop documents with the same index key as a previously indexed one
	Background       bool     // Build index in background and return immediately
	Sparse           bool     // Only index documents containing the Key fields
	ExpireAfter      time.Duration
	Name             string
	Min, Max         int
	Minf, Maxf       float64
	BucketSize       float64
	Bits             int
	DefaultLanguage  string
	LanguageOverride string
	Weights          map[string]int
	Collation        *mgo.Collation
}

// BulkResult See the BulkResult documentation in `gopkg.in/mgo.v2` for more information.
type BulkResult struct {
	Matched  int
	Modified int
}

func makeChangeInfo(in *mgo.ChangeInfo) *ChangeInfo {
	if in == nil {
		return nil
	}
	ret := ChangeInfo(*in)
	return &ret
}

func makeMapReduce(in *MapReduce) *mgo.MapReduce {
	if in == nil {
		return nil
	}
	ret := mgo.MapReduce(*in)
	return &ret
}

func makeMapReduceInfo(in *mgo.MapReduceInfo) *MapReduceInfo {
	if in == nil {
		return nil
	}
	ret := MapReduceInfo(*in)
	return &ret
}

func makeBulkResult(in *mgo.BulkResult) *BulkResult {
	if in == nil {
		return nil
	}

	return &BulkResult{
		Matched:  in.Matched,
		Modified: in.Modified,
	}
}
