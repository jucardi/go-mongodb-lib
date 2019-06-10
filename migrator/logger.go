package migrator

import "github.com/jucardi/go-mongodb-lib/log"

func SetLogger(logger log.ILogger) {
	log.Set(logger)
}
