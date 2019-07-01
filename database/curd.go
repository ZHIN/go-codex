package database

import (
	"database/sql"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
)

const getDBObjectErrorMessage = "无法获取数据库对象"

// Create 创建数据
func (r *DatabaseRepo) Create(item interface{}) error {
	r.put()
	defer r.pop()
	db, err := getDB(r.dbKey)

	if err != nil {
		return err
	}
	defer db.Close()

	db = db.Create(item)
	if err := db.Error; err != nil {
		return warpDBError(db, "DB.Create", fmt.Sprintf("DB.Create Error Conn:%s Type:%s", r.dbKey, reflect.TypeOf(item).String()))
	}
	return nil
}

// Save 保存数据
func (r *DatabaseRepo) Save(item interface{}) error {

	r.put()
	defer r.pop()

	db, err := getDB(r.dbKey)
	if err != nil {
		return err
	}
	defer db.Close()
	db = db.Save(item)
	if err := db.Error; err != nil {
		return warpDBError(db, "DB.Save", fmt.Sprintf("DB.Save Error Conn:%s Type:%s", r.dbKey, reflect.TypeOf(item).String()))
	}
	return nil
}

// Update 更新列
func (r *DatabaseRepo) Update(model interface{}, query string, params []interface{}, item ...interface{}) error {

	r.put()
	defer r.pop()

	db, err := getDB(r.dbKey)
	if err != nil {
		return err
	}

	defer db.Close()
	db = db.Model(model).Where(query, params...)
	db = db.Update(item...)
	if err := db.Error; err != nil {
		return warpDBError(db, "DB.Update", fmt.Sprintf("DB.Update Error Conn:%s Type:%s", r.dbKey, reflect.TypeOf(item).String()))
	}
	return nil

}

func (r *DatabaseRepo) UpdateColumn(model interface{}, query string, params []interface{}, attrs ...interface{}) error {

	r.put()
	defer r.pop()

	db, err := getDB(r.dbKey)
	if err != nil {
		return err
	}
	defer db.Close()

	db = db.Model(model).Where(query, params...).UpdateColumn(attrs...)
	if err := db.Error; err != nil {
		return warpDBError(db, "DB.UpdateColumn", fmt.Sprintf("DB.UpdateColumn Error Conn:%s Type:%s", r.dbKey, reflect.TypeOf(model).String()))
	}
	return nil
}

func (r *DatabaseRepo) Updates(model interface{}, query string, where []interface{}, item map[string]interface{}) error {

	r.put()
	defer r.pop()

	db, err := getDB(r.dbKey)
	if err != nil {
		return err
	}
	defer db.Close()

	db = db.Model(model).Where(query, where...).Updates(item)

	if err := db.Error; err != nil {
		return warpDBError(db, "DB.Updates", fmt.Sprintf("DB.Updates Error Conn:%s Type:%s", r.dbKey, reflect.TypeOf(item).String()))
	}
	return nil

}

func (r *DatabaseRepo) Delete(item interface{}, query string, params ...interface{}) error {

	r.put()
	defer r.pop()

	db, err := getDB(r.dbKey)
	if err != nil {
		return err
	}
	defer db.Close()
	query = strings.Trim(query, "")
	if query != "" {
		db = db.Where(query, params...)
	}
	db = db.Delete(item)
	if err := db.Error; err != nil {
		return warpDBError(db, "DB.Updates", fmt.Sprintf("DB.Delete Error Conn:%s Type:%s", r.dbKey, reflect.TypeOf(item).String()))
	}
	return nil
}

type SearchOption struct {
	Limit    int
	Offset   int
	Where    string
	Order    string
	Params   []interface{}
	TotalOut *int
}

type DBResult struct {
	Err              error
	IsRecordNotFound bool
}

func (r *DatabaseRepo) First(item interface{}, where string, params ...interface{}) DBResult {

	return r.FirstEX(item, SearchOption{
		Limit:  1,
		Offset: 0,
		Where:  where,
		Params: params,
	})
}
func (r *DatabaseRepo) FirstEX(item interface{}, option SearchOption) DBResult {
	r.put()
	defer r.pop()
	db, err := getDB(r.dbKey)
	if err != nil {
		return DBResult{Err: err}
	}
	defer db.Close()

	option.Order = strings.Trim(option.Order, "")
	if option.Order != "" {
		db = db.Order(option.Order)
	}
	values := []interface{}{}
	option.Where = strings.Trim(option.Where, "")
	if option.Where != "" {
		values = append(values, option.Where)
		values = append(values, option.Params...)
	}
	db = db.Offset(option.Offset)
	db = db.First(item, values...)
	if db.Error != nil {

		return DBResult{
			Err:              warpDBError(db, "DB.First", fmt.Sprintf("DB.First Error Conn:%s Type:%s WHERE:%v", r.dbKey, reflect.TypeOf(item).String(), option.Where)),
			IsRecordNotFound: false,
		}
	}

	return DBResult{IsRecordNotFound: db.RecordNotFound()}

}

func (r *DatabaseRepo) Find(list interface{}, offset int, limit int, where string, order string, params ...interface{}) error {
	return r.FindEX(list, SearchOption{
		Limit:  limit,
		Offset: offset,
		Where:  where,
		Params: params,
		Order:  order,
	})
}

func (r *DatabaseRepo) FindEX(list interface{}, option SearchOption) error {
	r.put()
	defer r.pop()
	db, err := getDB(r.dbKey)
	if err != nil {
		return err
	}
	defer db.Close()
	option.Order = strings.Trim(option.Order, "")

	values := []interface{}{}
	option.Where = strings.Trim(option.Where, "")
	if option.Where != "" {
		values = append(values, option.Where)
		values = append(values, option.Params...)
	}

	if option.Order != "" {
		db = db.Order(option.Order)
	}

	dbCount := db

	if option.TotalOut != nil {
		dbCount.Model(list).Where(option.Where, option.Params...).Count(option.TotalOut)
	}

	db = db.Offset(option.Offset)
	if option.Limit > 0 {
		db = db.Limit(option.Limit)
	} else {
		db = db.Limit(math.MaxInt32)
	}

	db = db.Find(list, values...)

	if db.Error != nil {
		return warpDBError(db, "DB.Find", fmt.Sprintf("DB.Find Error Conn:%s Type:%s WHERE:%v", r.dbKey, reflect.TypeOf(list).String(), option.Where))

	}
	return nil

}

func (r *DatabaseRepo) Count(item interface{}, query string, values ...interface{}) (int, *DbError) {

	r.put()
	defer r.pop()

	db, err := getDB(r.dbKey)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	total := 0

	db = db.Model(item).Where(query, values...).Count(&total)
	if err := db.Error; err != nil {
		return -1, warpDBError(db, "DB.Count", fmt.Sprintf("DB.Count Error Conn:%s SQL:%s WHERE:%s... ", r.dbKey, query, values))
	}
	return total, nil
}

func (r *DatabaseRepo) Exec(sql string, values ...interface{}) error {

	r.put()
	defer r.pop()

	db, err := getDB(r.dbKey)
	if err != nil {
		return err
	}
	defer db.Close()
	db = db.Exec(sql, values...)
	if db.Error != nil {
		return warpDBError(db, "DB.Exec", fmt.Sprintf("DB.Exec Error Conn:%s SQL:%s WHERE:%s...", r.dbKey, sql, values))
	}
	return nil
}

type RowScanHandler func(*gorm.DB, *sql.Rows) error
type TrasnsationInvokeHandler func(db *gorm.DB) error

func (r *DatabaseRepo) InvokeTransation(callback TrasnsationInvokeHandler) error {
	r.put()
	defer r.pop()
	db, err := getDB(r.dbKey)

	if err != nil {
		return err
	}
	defer db.Close()
	var e = callback(db)
	if e != nil {
		return warpDBError(db, "DB.InvokeTransation", e.Error())
	}
	return nil
}

func (r *DatabaseRepo) RawSelect(rawSQL string, rowScanCallback RowScanHandler, values ...interface{}) error {

	r.put()
	defer r.pop()

	db, err2 := getDB(r.dbKey)

	if err2 != nil {
		return err2
	}

	defer db.Close()
	var rows *sql.Rows
	var err error

	if rows, err = db.Raw(rawSQL, values...).Rows(); err != nil {
		return warpDBError(db, "DB.RawSelect", fmt.Sprintf("DB.RawSelect Error Conn:%s SQL:%s WHERE:%s...", r.dbKey, rawSQL, values))
	}

	defer rows.Close()
	for rows.Next() {
		if rowScanCallback != nil {
			err = rowScanCallback(db, rows)
			if err != nil {
				return warpDBError(db, "DB.RawSelect", err.Error())

			}
		}
	}
	return nil
}

func (r *DatabaseRepo) ExecuteScalar(rawSQL string, params []interface{}, values ...interface{}) error {

	r.put()
	defer r.pop()
	db, err2 := getDB(r.dbKey)

	if err2 != nil {
		return err2
	}

	defer db.Close()
	var rows *sql.Rows
	var err error
	if rows, err = db.Raw(rawSQL, params...).Rows(); err != nil {
		return warpDBError(db, "DB.RawSelect", fmt.Sprintf("DB.ExecuteScalar Error Conn:%s SQL:%s WHERE:%s...", r.dbKey, rawSQL, params))
	}

	defer rows.Close()
	for rows.Next() {
		rows.Scan(values...)
		break
	}
	return nil
}
