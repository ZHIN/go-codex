package database

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
)

var dbKeyPoool = map[string]DBSet{}

type DBSet struct {
	DBType             string
	DBConnectionString string
	MaxOpenConns       int
	MaxIdleConns       int
}

func SetDB(dbKey string, _dbType string, connStr string) {

	dbKeyPoool[dbKey] = DBSet{
		DBType:             _dbType,
		DBConnectionString: connStr,
	}
}

func SetDBSet(dbKey string, opt DBSet) {
	dbKeyPoool[dbKey] = opt
}

func GetDB(dbKey string) (*gorm.DB, error) {

	if set, found := dbKeyPoool[dbKey]; found {
		db, err := gorm.Open(set.DBType, set.DBConnectionString)
		if err != nil {
			return nil, err
		}
		db.DB().SetMaxOpenConns(set.MaxOpenConns)
		db.DB().SetMaxIdleConns(set.MaxIdleConns)
		return db, nil
	}

	return nil, fmt.Errorf("找不到数据库相关连接配置（%s）", dbKey)

}

func GetDBNoErr(dbKey string) *gorm.DB {
	db, err := GetDB(dbKey)
	if err != nil {
		panic(err)
	}
	return db
}

func AutoMigrate(dbKey string, values ...interface{}) {
	db, err := GetDB(dbKey)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = db.AutoMigrate(values...).Error; err != nil {
		for _, value := range values {
			logrus.Warnf("AutoMigrate Error (%s)  Type:%v", err.Error(), reflect.TypeOf(value))
		}
	}
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
