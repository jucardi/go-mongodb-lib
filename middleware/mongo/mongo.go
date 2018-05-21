package mongo

import (
	"github.com/gin-gonic/gin"
	"github.com/jucardi/go-mongodb-lib/log"
	"github.com/jucardi/go-mongodb-lib/mgo"
	"github.com/sirupsen/logrus"
	"github.com/jucardi/go-mongodb-lib/migrator"
)

const (
	mongoContextStore = "mongoStore"
)

type context struct {
	session  mgo.ISession
	database string
}

// DBStore is the structure that holds the database instance used in the Gin context.
type DBStore struct {
	db mgo.IDatabase
}

var currentContext *context

// Use receives an implementation of the gig.IRoutes (normally a gin engine) and adds the mongodb middleware to be used for
// every context created by the router.
//
// {url} - the url connection string for the mongodb instance.
func Use(router gin.IRoutes, url string) gin.IRoutes {
	return router.Use(CreateHandler(url))
}

// // UseWithMigration like UseMongo, but enables reading the migration directory for migration scripts,
// // and ensures the DB is up to date with these scripts. If the contents on one previously migrated script
// // is changed, it will panic because of the integrity of the db would have been compromised.
// //
// // {mongoUrl}            - the mongo connection url.
// // {migrationDir}        - the directory where the migration scripts are stored.
// // {failOnOrderMismatch} - guarantees the alphabetical order of script execution. If a new script is added
// //                         and by alphabetical order falls before a script that has been previously migrated
// //                         the migration will fail if this flag is set to 'true'
func UseWithMigration(router gin.IRoutes, mongoUrl, migrationDir string, failOnOrderMismatch bool) gin.IRoutes {
	h := Use(router, mongoUrl)
	s, db := GetDb()

	defer s.Close()
	if err := migrator.Migrate(migrationDir, db, failOnOrderMismatch); err != nil {
		logrus.Fatalf("An error occurred while migrating data. %s", err.Error())
	}

	return h
}

// CreateHandler associates the specified mongo url which is used to connect to MongoDB
// A gin handler is returned which can be applied to middleware to start and end sessions
// to mongo on every request.
//
// {url} - the url connection string for the mongodb instance.
func CreateHandler(url string) gin.HandlerFunc {
	s, err := mgo.Dial(url)

	if err != nil {
		log.Get().Panic("Unable to establish a connection with mongo. Max retries reached.")
	}

	s.SetDefaultSafe()

	dbCtx := &context{session: s, database: s.DB("").Name()}
	currentContext = dbCtx

	return dbCtx.handler
}

// Session returns a pointer to the global Mongo session
func Session() mgo.ISession {
	return currentContext.session
}

// DatabaseName returns the name of the database stored in the context
func DatabaseName() string {
	return currentContext.database
}

// GetDb gets a clone of the DB session generated for the current gin.Context, wrapped in the new implementation of the go-foundation mgo package.
//
//    Usage:
//            s, db := GetDb()
//            defer s.Close()
//
//    NOTE: This function does not close the clone session. Must be manually closed.
//
func GetDb() (mgo.ISession, mgo.IDatabase) {
	session := Session().Clone()
	db := session.DB("")
	return session, db
}

func (d *context) handler(c *gin.Context) {
	c.Next()
	if db, exists := c.Get(mongoContextStore); exists {
		db.(*DBStore).db.Session().Close()
	}
}

// Get returns the dbStore instance where the db information is held
func Get(c *gin.Context) *DBStore {
	if db, exists := c.Get(mongoContextStore); exists {
		return db.(*DBStore)
	}

	s := Session().Clone()
	db := &DBStore{s.DB("")}
	c.Set(mongoContextStore, db)
	return db
}

// DB returns a database from the current session which can be used
// for queries.  It will automatically be closed when the current request
// completes by the middleware handler
func (d *DBStore) DB() mgo.IDatabase {
	return d.db
}

// IsErrNotFound indicates if the error was generated due to a record not been found.
func IsErrNotFound(err error) bool {
	return err.Error() == mgo.ErrNotFound.Error()
}
