package migrator

import (
	"errors"
	"testing"

	"github.com/jucardi/go-mongodb-lib/mgo"
	"github.com/jucardi/go-mongodb-lib/testutils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

const migrationPath = "./test_assets/db_migration"

func TestMigrateSuccess(t *testing.T) {
	db := mgo.MockDb()
	col := mgo.MockCollection()
	q := mgo.MockQuery()

	db.WhenC(MigrationCollection, col)
	col.WhenFind(bson.M{}, q)

	err := Migrate(migrationPath, db, true)
	assert.Equal(t, 1, col.Times("Find"))
	assert.Equal(t, 2, col.Times("Insert"))
	assert.Equal(t, 3, db.Times("C"))
	assert.Equal(t, 2, db.Times("Run"))
	assert.Equal(t, 1, q.Times("Sort"))
	assert.Equal(t, 1, q.Times("All"))

	assert.Nil(t, err)
}

func TestReadMigrateCollectionFailed(t *testing.T) {
	msg := "some error"
	db := mgo.MockDb()
	col := mgo.MockCollection()
	q := mgo.MockQuery()

	db.WhenC(MigrationCollection, col)
	col.WhenFind(bson.M{}, q)
	q.WhenReturn("All", errors.New(msg))

	err := Migrate(migrationPath, db, true)
	assert.NotNil(t, err)
	assert.True(t, err.Is(ErrDbOperation))
	assert.Equal(t, "Unable to read Database info. "+msg, err.Error())

	assert.Equal(t, 1, col.Times("Find"))
	assert.Equal(t, 0, col.Times("Insert"))
	assert.Equal(t, 1, db.Times("C"))
	assert.Equal(t, 0, db.Times("Run"))
	assert.Equal(t, 1, q.Times("Sort"))
	assert.Equal(t, 1, q.Times("All"))
}

func TestReadScriptPathFailed(t *testing.T) {
	path := "some-invalid-path"
	db := mgo.MockDb()
	col := mgo.MockCollection()
	q := mgo.MockQuery()

	db.WhenC(MigrationCollection, col)
	col.WhenFind(bson.M{}, q)

	err := Migrate(path, db, true)
	assert.NotNil(t, err)
	assert.True(t, err.Is(ErrFileAccess))
	assert.Equal(t, "Unable to access scripts path. open "+path+": no such file or directory", err.Error())

	assert.Equal(t, 1, col.Times("Find"))
	assert.Equal(t, 0, col.Times("Insert"))
	assert.Equal(t, 1, db.Times("C"))
	assert.Equal(t, 0, db.Times("Run"))
	assert.Equal(t, 1, q.Times("Sort"))
	assert.Equal(t, 1, q.Times("All"))
}

func TestPreviousDataSuccess(t *testing.T) {
	db := mgo.MockDb()
	col := mgo.MockCollection()
	q := mgo.MockQuery()

	db.WhenC(MigrationCollection, col)
	col.WhenFind(bson.M{}, q)

	q.When("All", func(args ...interface{}) []interface{} {
		result := args[0]
		list := result.(*[]*MigrationInfo)
		*list = append(*list, &MigrationInfo{
			ScriptId: "script_001.js",
			Hash:     "b280f134425a4153026cf227069d4cc1",
		})
		return testutils.MakeReturn(nil)
	})

	assert.Nil(t, Migrate(migrationPath, db, true))

	assert.Equal(t, 1, col.Times("Find"))
	assert.Equal(t, 1, col.Times("Insert"))
	assert.Equal(t, 2, db.Times("C"))
	assert.Equal(t, 1, db.Times("Run"))
	assert.Equal(t, 1, q.Times("Sort"))
	assert.Equal(t, 1, q.Times("All"))
}

func TestPreviousDataFailedHash(t *testing.T) {
	db := mgo.MockDb()
	col := mgo.MockCollection()
	q := mgo.MockQuery()

	db.WhenC(MigrationCollection, col)
	col.WhenFind(bson.M{}, q)
	q.When("All", func(args ...interface{}) []interface{} {
		result := args[0]
		list := result.(*[]*MigrationInfo)
		*list = append(*list, &MigrationInfo{
			ScriptId: "script_001.js",
			Hash:     "1234",
		})
		return testutils.MakeReturn(nil)
	})

	err := Migrate(migrationPath, db, true)
	assert.NotNil(t, err)
	assert.Equal(t, "File 'script_001.js' was previously migrated but hashes don't match.", err.Error())

	assert.Equal(t, 1, col.Times("Find"))
	assert.Equal(t, 0, col.Times("Insert"))
	assert.Equal(t, 1, db.Times("C"))
	assert.Equal(t, 0, db.Times("Run"))
	assert.Equal(t, 1, q.Times("Sort"))
	assert.Equal(t, 1, q.Times("All"))
}

func TestRunFailed(t *testing.T) {
	db := mgo.MockDb()
	col := mgo.MockCollection()
	q := mgo.MockQuery()

	db.WhenC(MigrationCollection, col)
	col.WhenFind(bson.M{}, q)
	db.WhenRun(nil, errors.New("some error"))

	err := Migrate(migrationPath, db, true)
	assert.NotNil(t, err)
	assert.Equal(t, "Unable to run command 'script_001.js'. some error", err.Error())

	assert.Equal(t, 1, col.Times("Find"))
	assert.Equal(t, 0, col.Times("Insert"))
	assert.Equal(t, 1, db.Times("C"))
	assert.Equal(t, 1, db.Times("Run"))
	assert.Equal(t, 1, q.Times("Sort"))
	assert.Equal(t, 1, q.Times("All"))
}

func TestInsertFailed(t *testing.T) {
	db := mgo.MockDb()
	col := mgo.MockCollection()
	q := mgo.MockQuery()

	db.WhenC(MigrationCollection, col)
	col.WhenFind(bson.M{}, q)
	col.WhenInsert(nil, errors.New("some error"))

	err := Migrate(migrationPath, db, true)
	assert.NotNil(t, err)
	assert.Equal(t, "Unable to save migration info for 'script_001.js'", err.Error())

	assert.Equal(t, 1, col.Times("Find"))
	assert.Equal(t, 1, col.Times("Insert"))
	assert.Equal(t, 2, db.Times("C"))
	assert.Equal(t, 1, db.Times("Run"))
	assert.Equal(t, 1, q.Times("Sort"))
	assert.Equal(t, 1, q.Times("All"))
}
