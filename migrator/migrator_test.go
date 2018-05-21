package migrator

import (
	"github.com/jucardi/go-mongodb-lib/mgo"
	"github.com/jucardi/go-mongodb-lib/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

const migrationPath = "./test_assets/db_migration"

func TestMigrateSuccess(t *testing.T) {
	db := mgo.MockDb(t)
	col := mgo.MockCollection(t)
	q := mgo.MockQuery(t)

	db.WhenC(MigrationCollection, col)
	col.WhenFind(bson.M{}, q)

	err := Migrate(migrationPath, db, true)
	col.Times("Find", 1)
	col.Times("Insert", 2)
	db.Times("C", 3)
	db.Times("Run", 2)
	q.Times("Sort", 1)
	q.Times("All", 1)

	assert.Nil(t, err)
}

func TestReadMigrateCollectionFailed(t *testing.T) {
	msg := "some error"
	db := mgo.MockDb(t)
	col := mgo.MockCollection(t)
	q := mgo.MockQuery(t)

	db.WhenC(MigrationCollection, col)
	col.WhenFind(bson.M{}, q)
	q.When("All", func(t *testing.T, args ...interface{}) []interface{} {
		return testutils.MakeReturn(errors.New(msg))
	})
	err := Migrate(migrationPath, db, true)
	assert.NotNil(t, err)

	assert.True(t, err.Is(ErrDbOperation))
	assert.Equal(t, "Unable to read Database info. "+msg, err.Error())
	col.Times("Find", 1)
	col.Times("Insert", 0)
	db.Times("C", 1)
	db.Times("Run", 0)
	q.Times("Sort", 1)
	q.Times("All", 1)
}

func TestReadScriptPathFailed(t *testing.T) {
	path := "some-invalid-path"
	db := mgo.MockDb(t)
	col := mgo.MockCollection(t)
	q := mgo.MockQuery(t)

	db.WhenC(MigrationCollection, col)
	col.WhenFind(bson.M{}, q)

	err := Migrate(path, db, true)

	assert.True(t, err.Is(ErrFileAccess))
	assert.Equal(t, "Unable to access scripts path. open "+path+": no such file or directory", err.Error())
	col.Times("Find", 1)
	col.Times("Insert", 0)
	db.Times("C", 1)
	db.Times("Run", 0)
	q.Times("Sort", 1)
	q.Times("All", 1)
}

func TestPreviousDataSuccess(t *testing.T) {
	db := mgo.MockDb(t)
	col := mgo.MockCollection(t)
	q := mgo.MockQuery(t)

	db.WhenC(MigrationCollection, col)
	col.WhenFind(bson.M{}, q)

	q.When("All", func(t *testing.T, args ...interface{}) []interface{} {
		result := args[0]
		list := result.(*[]*MigrationInfo)
		*list = append(*list, &MigrationInfo{
			ScriptId: "script_001.js",
			Hash:     "b280f134425a4153026cf227069d4cc1",
		})
		return testutils.MakeReturn(nil)
	})

	err := Migrate(migrationPath, db, true)
	col.Times("Find", 1)
	col.Times("Insert", 1)
	db.Times("C", 2)
	db.Times("Run", 1)
	q.Times("Sort", 1)
	q.Times("All", 1)
	assert.Nil(t, err)
}

func TestPreviousDataFailedHash(t *testing.T) {
	db := mgo.MockDb(t)
	col := mgo.MockCollection(t)
	q := mgo.MockQuery(t)

	db.WhenC(MigrationCollection, col)
	col.WhenFind(bson.M{}, q)
	q.When("All", func(t *testing.T, args ...interface{}) []interface{} {
		result := args[0]
		list := result.(*[]*MigrationInfo)
		*list = append(*list, &MigrationInfo{
			ScriptId: "script_001.js",
			Hash:     "1234",
		})
		return testutils.MakeReturn(nil)
	})

	err := Migrate(migrationPath, db, true)
	col.Times("Find", 1)
	col.Times("Insert", 0)
	db.Times("C", 1)
	db.Times("Run", 0)
	q.Times("Sort", 1)
	q.Times("All", 1)
	assert.NotNil(t, err)
	assert.Equal(t, "File 'script_001.js' was previously migrated but hashes don't match.", err.Error())
}

func TestRunFailed(t *testing.T) {
	db := mgo.MockDb(t)
	col := mgo.MockCollection(t)
	q := mgo.MockQuery(t)

	db.WhenC(MigrationCollection, col)
	col.WhenFind(bson.M{}, q)
	db.WhenRun(nil, errors.New("some error"))

	err := Migrate(migrationPath, db, true)
	col.Times("Find", 1)
	col.Times("Insert", 0)
	db.Times("C", 1)
	db.Times("Run", 1)
	q.Times("Sort", 1)
	q.Times("All", 1)
	assert.NotNil(t, err)
	assert.Equal(t, "Unable to run command 'script_001.js'. some error", err.Error())
}

func TestInsertFailed(t *testing.T) {
	db := mgo.MockDb(t)
	col := mgo.MockCollection(t)
	q := mgo.MockQuery(t)

	db.WhenC(MigrationCollection, col)
	col.WhenFind(bson.M{}, q)
	col.WhenInsert(nil, errors.New("some error"))

	err := Migrate(migrationPath, db, true)
	col.Times("Find", 1)
	col.Times("Insert", 1)
	db.Times("C", 2)
	db.Times("Run", 1)
	q.Times("Sort", 1)
	q.Times("All", 1)
	assert.NotNil(t, err)
	assert.Equal(t, "Unable to save migration info for 'script_001.js'", err.Error())
}
