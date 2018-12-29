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
func (r *DatabaseRepo) Create(item interface{}) *DbError {
	r.put()
	defer r.pop()
	db, err := getDB(r.dbKey)

	if err != nil {
		return err
	}
	defer db.Close()
	if err := db.Create(item).Error; err != nil {
		return warpDBError(db, "DB.Create", fmt.Sprintf("DB.Create Error Conn:%s Type:%s", r.dbKey, reflect.TypeOf(item).String()))
	}
	return warpDBError(db, "DB.Create", "")
}

// Save 保存数据
func (r *DatabaseRepo) Save(item interface{}) *DbError {

	r.put()
	defer r.pop()

	db, err := getDB(r.dbKey)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := db.Save(item).Error; err != nil {
		return warpDBError(db, "DB.Save", fmt.Sprintf("DB.Save Error Conn:%s Type:%s", r.dbKey, reflect.TypeOf(item).String()))
	}
	return warpDBError(db, "DB.Save", "")
}

// Update 更新列
func (r *DatabaseRepo) Update(model interface{}, query string, params []interface{}, item ...interface{}) *DbError {

	r.put()
	defer r.pop()

	db, err := getDB(r.dbKey)
	if err != nil {
		return err
	}

	defer db.Close()
	db = db.Model(model).Where(query, params...)
	if err := db.Update(item...).Error; err != nil {
		return warpDBError(db, "DB.Update", fmt.Sprintf("DB.Update Error Conn:%s Type:%s", r.dbKey, reflect.TypeOf(item).String()))
	}
	return warpDBError(db, "DB.Update", "")

}

func (r *DatabaseRepo) UpdateColumn(model interface{}, query string, params []interface{}, attrs ...interface{}) *DbError {

	r.put()
	defer r.pop()

	db, err := getDB(r.dbKey)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := db.Model(model).Where(query, params...).UpdateColumn(attrs...).Error; err != nil {
		return warpDBError(db, "DB.UpdateColumn", fmt.Sprintf("DB.UpdateColumn Error Conn:%s Type:%s", r.dbKey, reflect.TypeOf(model).String()))
	}
	return warpDBError(db, "DB.UpdateColumn", "")
}

func (r *DatabaseRepo) Updates(model interface{}, query string, where []interface{}, item map[string]interface{}) *DbError {

	r.put()
	defer r.pop()

	db, err := getDB(r.dbKey)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := db.Model(model).Where(query, where...).Updates(item).Error; err != nil {
		return warpDBError(db, "DB.Updates", fmt.Sprintf("DB.Updates Error Conn:%s Type:%s", r.dbKey, reflect.TypeOf(item).String()))
	}
	return warpDBError(db, "DB.Updates", "")

}

func (r *DatabaseRepo) Delete(item interface{}, query string, params ...interface{}) *DbError {

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
	if err := db.Delete(item).Error; err != nil {
		return warpDBError(db, "DB.Updates", fmt.Sprintf("DB.Delete Error Conn:%s Type:%s", r.dbKey, reflect.TypeOf(item).String()))
	}
	return warpDBError(db, "DB.Delete", "")
}

type SearchOption struct {
	Limit    int
	Offset   int
	Where    string
	Order    string
	Params   []interface{}
	TotalOut *int
}

func (r *DatabaseRepo) First(item interface{}, option SearchOption) *DbError {
	r.put()
	defer r.pop()
	db, err := getDB(r.dbKey)
	if err != nil {
		return err
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

	if db.Error != nil && !db.RecordNotFound() {
		return warpDBError(db, "DB.First", fmt.Sprintf("DB.First Error Conn:%s Type:%s WHERE:%v", r.dbKey, reflect.TypeOf(item).String(), option.Where))
	}
	return warpDBError(db, "DB.First", "")
}

func (r *DatabaseRepo) Find(item interface{}, option SearchOption) *DbError {
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
		dbCount.Model(item).Where(option.Where, option.Params).Count(option.TotalOut)
	}

	db = db.Offset(option.Offset)
	if option.Limit > 0 {
		db = db.Limit(option.Limit)
	} else {
		db = db.Limit(math.MaxInt32)
	}

	db = db.Find(item, values...)

	if db.Error != nil && !db.RecordNotFound() {
		return warpDBError(db, "DB.Find", fmt.Sprintf("DB.Find Error Conn:%s Type:%s WHERE:%v", r.dbKey, reflect.TypeOf(item).String(), option.Where))

	}

	return warpDBError(db, "DB.Find", "")

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
	if err := db.Model(item).Where(query, values...).Count(&total).Error; err != nil {
		return -1, warpDBError(db, "DB.Count", fmt.Sprintf("DB.Count Error Conn:%s SQL:%s WHERE:%s... ", r.dbKey, query, values))
	}
	return total, warpDBError(db, "DB.Count", "")
}

func (r *DatabaseRepo) Exec(sql string, values ...interface{}) *DbError {

	r.put()
	defer r.pop()

	db, err := getDB(r.dbKey)
	if err != nil {
		return err
	}
	defer db.Close()
	db = db.Exec(sql, values...)
	return warpDBError(db, "DB.Exec", fmt.Sprintf("DB.Exec Error Conn:%s SQL:%s WHERE:%s...", r.dbKey, sql, values))
}

type RowScanHandler func(*gorm.DB, *sql.Rows) error
type TrasnsationInvokeHandler func(db *gorm.DB) error

func (r *DatabaseRepo) InvokeTransation(callback TrasnsationInvokeHandler) *DbError {
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
	return warpDBError(db, "DB.InvokeTransation", "")
}

func (r *DatabaseRepo) RawSelect(rawSQL string, rowScanCallback RowScanHandler, values ...interface{}) *DbError {

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
	return warpDBError(db, "DB.RawSelect", "")
}

func (r *DatabaseRepo) ExecuteScalar(rawSQL string, values ...interface{}) *DbError {

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
		rows.Scan(values...)
		break
	}
	return warpDBError(db, "DB.RawSelect", "")
}
