package mgo

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/jucardi/go-mongodb-lib/log"
	"github.com/jucardi/go-mongodb-lib/mgo"
	"github.com/jucardi/go-osx/paths"
	"github.com/jucardi/go-streams/streams"
	"gopkg.in/mgo.v2/bson"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	MigrationCollection = "_migration"
	ErrFileAccess       = 0x1
	ErrDbAccess         = 0x2
	ErrDbOperation      = 0x4
	ErrOrderFailed      = 0x8
	ErrHashingFailed    = 0x16
)

type MigrationErrorCode int

type MigrationError struct {
	Message string
	Code    int
}

// Returns the error message
func (m *MigrationError) Error() string {
	return m.Message
}

func (m *MigrationError) ErrorCode() int {
	return m.Code
}

func (m *MigrationError) Is(errorCode int) bool {
	return m.Code&errorCode == errorCode
}

type MigrationInfo struct {
	ScriptId  string    `json:"script_id" bson:"script_id"`
	Hash      string    `json:"hash" bson:"hash"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

func Migrate(dataDir string, db mgo.IDatabase, failOnOrderMismatch bool) *MigrationError {
	var infos []*MigrationInfo

	if err := db.C(MigrationCollection).Find(bson.M{}).Sort("filename").All(&infos); err != nil {
		return &MigrationError{
			Message: fmt.Sprintf("Unable to read Database info. %s", err.Error()),
			Code:    ErrDbAccess | ErrDbOperation,
		}
	}

	objs, err := ioutil.ReadDir(dataDir)

	if err != nil {
		return &MigrationError{
			Message: fmt.Sprintf("Unable to access scripts path. %s", err.Error()),
			Code:    ErrFileAccess,
		}
	}

	foundNonMigrated := false

	var toMigrate []MigrationInfo

	for _, f := range streams.From(objs).OrderBy(func(a, b interface{}) int {
		return strings.Compare(b.(os.FileInfo).Name(), a.(os.FileInfo).Name())
	}, true).ToArray().([]os.FileInfo) {

		if f.IsDir() {
			continue
		}

		log.Get().Infof("Migrating file '%s'", f.Name())
		fullPath := paths.Combine(dataDir, f.Name())
		hash, hashErr := computeHash(fullPath)

		if hashErr != nil {
			return &MigrationError{
				Message: fmt.Sprintf("Error computing hash for file '%s', aborting migration.", hashErr.Error()),
				Code:    ErrHashingFailed | ErrFileAccess,
			}
		}

		// Contains the migration info
		if inf := streams.From(infos).
			Filter(
				func(obj interface{}) bool {
					return obj.(*MigrationInfo).ScriptId == f.Name()
				}).
			First(); inf != nil {

			if foundNonMigrated && failOnOrderMismatch {
				return &MigrationError{
					Message: fmt.Sprintf("Non-Migrated file found before '%s' which has been migrated. Order import failed, unable to proceed.", f.Name()),
					Code:    ErrOrderFailed,
				}
			}

			info := inf.(*MigrationInfo)

			if info.Hash != hash {
				return &MigrationError{
					Message: fmt.Sprintf("File '%s' was previously migrated but hashes don't match.", f.Name()),
					Code:    ErrHashingFailed,
				}
			} else {
				log.Get().Infof("File '%s' previously migrated, continuing", f.Name())
			}

			continue

		} else {
			foundNonMigrated = true
			toMigrate = append(toMigrate, MigrationInfo{
				ScriptId: f.Name(),
				Hash:     hash,
			})
		}
	}

	for _, info := range toMigrate {
		fullPath := paths.Combine(dataDir, info.ScriptId)
		if content, err := ioutil.ReadFile(fullPath); err != nil {
			return &MigrationError{
				Message: fmt.Sprintf("Unable to read data file '%s': %s", info.ScriptId, err.Error()),
				Code:    ErrFileAccess,
			}
		} else {

			jsContent := string(content)
			var resp map[string]interface{}

			if err := db.Run(bson.M{"eval": jsContent}, &resp); err != nil {
				return &MigrationError{
					Message: fmt.Sprintf("Unable to run command '%s'. %s", info.ScriptId, err.Error()),
					Code:    ErrDbOperation,
				}
			}

			log.Get().Debug(resp)
			info.Timestamp = time.Now()

			if err := db.C(MigrationCollection).Insert(info); err != nil {
				return &MigrationError{
					Message: fmt.Sprintf("Unable to save migration info for '%s'", info.ScriptId),
					Code:    ErrDbAccess | ErrDbOperation,
				}
			}
		}
	}

	return nil
}

func computeHash(filePath string) (string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return "", err
	}

	defer file.Close()
	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)[:16]
	return hex.EncodeToString(hashInBytes), nil
}
