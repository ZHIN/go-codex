package database

import (
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
)

var dbKeyPoool = map[string]DBSetOption{}

type DBSetOption struct {
	DBType             string
	DBConnectionString string
	MaxOpenConns       int
	MaxIdleConns       int
}

func SetDB(dbKey string, _dbType string, connStr string) {

	dbKeyPoool[dbKey] = DBSetOption{
		DBType:             _dbType,
		DBConnectionString: connStr,
	}
}

func SetDBSet(dbKey string, opt DBSetOption) {
	dbKeyPoool[dbKey] = opt
}

func getDB(dbKey string) (*gorm.DB, *DbError) {

	if set, found := dbKeyPoool[dbKey]; found {
		db, err := gorm.Open(set.DBType, set.DBConnectionString)
		if err != nil {
			return nil, warpDBError(db, "DB.OPEN", "")
		}
		db.DB().SetMaxOpenConns(set.MaxOpenConns)
		db.DB().SetMaxIdleConns(set.MaxIdleConns)
		return db, nil
	}
	return nil, warpDBError(nil, "DB.OPEN", fmt.Sprintf("找不到数据库相关连接配置（%s）", dbKey))

}

func AutoMigrate(dbKey string, values ...interface{}) *DbError {
	db, err := getDB(dbKey)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.AutoMigrate(values...).Error; err != nil {
		return warpDBError(db, "DB.AutoMigrate", "")
	}
	return nil
}

type DatabaseRepo struct {
	dbKey   string
	channel chan int
}

func (r *DatabaseRepo) put() {
	r.channel <- 0
}

func (r *DatabaseRepo) pop() {
	<-r.channel
}

var repos = map[string]*DatabaseRepo{}
var dblock = sync.Mutex{}

func Choice(dbKey string) *DatabaseRepo {

	dblock.Lock()
	defer dblock.Unlock()
	if item, found := repos[dbKey]; found {
		return item
	}

	maxOpenNum := 1
	if dbKeyPoool[dbKey].MaxOpenConns > 0 {
		maxOpenNum = dbKeyPoool[dbKey].MaxOpenConns
	}
	repos[dbKey] = &DatabaseRepo{dbKey: dbKey,
		channel: make(chan int, maxOpenNum),
	}
	return repos[dbKey]
}

type DbError struct {
	db      *gorm.DB
	tag     string
	message string
}

func (s *DbError) RecordNotFound() bool {
	return s.db.RecordNotFound()
}

func (s *DbError) Error() string {
	if s.db.Error == nil {
		return fmt.Sprintf("database error tag:%s message:%s", s.tag, s.message)
	}
	return fmt.Sprintf("database error:%s tag:%s message:%s", s.db.Error.Error(), s.tag, s.message)
}

func warpDBError(db *gorm.DB, tag string, message string) *DbError {

	if db == nil {
		return &DbError{db: db, tag: tag, message: message}
	}
	if db.Error != nil {
		var err = &DbError{db: db, tag: tag, message: message}
		if !db.RecordNotFound() {
			triggerErrorHandles(err)
		}
		return err
	}
	return nil
}

type DBErrorHandle func(err *DbError)

var errorHandles = []DBErrorHandle{}

func SetDBErrorHook(handle DBErrorHandle) {
	errorHandles = append(errorHandles, handle)
}

func triggerErrorHandles(err *DbError) {
	if errorHandles != nil {
		for _, handle := range errorHandles {
			if handle != nil {
				handle(err)
			}
		}
	}
}
