package database

import (
	"sync"
	"testing"
	"time"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Address struct {
	ID      int    `gorm:"column:id;primary_key;auto_increment"`
	Content string `gorm:"column:content"`
}

const defaultDBKey = "DEFAULT"

func setupTestCase(t *testing.T) {

	SetDBSet(defaultDBKey, DBSetOption{
		DBType:             "sqlite3",
		DBConnectionString: "data.db",
		MaxOpenConns:       1,
		MaxIdleConns:       1,
	})

	AutoMigrate(defaultDBKey, &Address{})

}

func TestDatabaseRepo_Create(t *testing.T) {
	setupTestCase(t)
	var err *DbError

	err = Choice(defaultDBKey).Create(&Address{Content: time.Now().String()})

	if err != nil {
		t.Error(err)
	}
}

func TestDatabaseRepo_CreateMany(t *testing.T) {
	setupTestCase(t)
	var err *DbError

	err = Choice(defaultDBKey).InvokeTransation(func(db *gorm.DB) error {
		tx := db.Begin()
		cnt := 10
		for cnt > 0 {
			cnt--
			if err := tx.Create(&Address{Content: time.Now().String()}).Error; err != nil {
				tx.Rollback()
				return err
			}

		}
		tx.Commit()
		return nil
	})

	if err != nil {
		t.Error(err)
	}
}

func TestDatabaseRepo_Find(t *testing.T) {
	setupTestCase(t)
	var err *DbError

	var list []Address
	var total int

	err = Choice(defaultDBKey).Find(&list, SearchOption{
		Offset:   2,
		Limit:    -1,
		TotalOut: &total,
		Where:    "id < 10",
		Order:    "id DESC",
	})

	if err != nil {
		t.Error(err)
	}

	t.Log("list:", list)
	t.Log("total:", total)
}

func TestDatabaseRepo_First(t *testing.T) {
	setupTestCase(t)
	var err *DbError

	var item Address

	err = Choice(defaultDBKey).First(&item, SearchOption{
		Offset: 2,
		Limit:  3,
		Where:  "id < 1",
		Order:  "id DESC",
	})

	if err != nil && !err.RecordNotFound() {
		t.Error(err)
	}
	t.Log("item:", item)
}

func TestDatabaseRepo_ErrorHook(t *testing.T) {
	setupTestCase(t)

	SetDBErrorHook(func(err *DbError) {
		t.Log("DBHOOK-ERR", err)
	})

	err := Choice(defaultDBKey).Exec("SELECT * FROM SDFDF")
	if err == nil {
		t.Error("err should not nil")
	}
}

func TestDatabaseRepo_CreateManyCronCurrent(t *testing.T) {
	setupTestCase(t)
	num := 500
	var wg sync.WaitGroup
	for num > 0 {
		num--
		wg.Add(1)
		go func(_wg *sync.WaitGroup) {
			cnt := 100
			for cnt > 0 {
				cnt--
				if err := Choice(defaultDBKey).Create(&Address{Content: time.Now().String()}); err != nil {
					t.Error(err)
					return
				}
			}
			_wg.Done()
		}(&wg)
	}
	wg.Wait()

}
