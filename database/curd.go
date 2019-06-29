package database

import (
	"database/sql"
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/zhin/go-codex/cerror"
)

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
func (r *DatabaseRepo) Create(item interface{}) error {

	r.put()
	defer r.pop()

	db, err := GetDB(r.dbKey)
	if err != nil {
		return cerror.NewCodeError(cerror.DB_ERROR, err)
	}
	defer db.Close()
	if err := db.Create(item).Error; err != nil {
		logrus.Warnf("DB.Create Error Conn:%s Type:%s (%s)", r.dbKey, reflect.TypeOf(item).String(), err.Error())
		return cerror.NewCodeError(cerror.DB_ERROR, err)
	}
	return nil
}

func (r *DatabaseRepo) Save(item interface{}) error {

	r.put()
	defer r.pop()

	db, err := GetDB(r.dbKey)
	if err != nil {
		return cerror.NewCodeError(cerror.DB_ERROR, err)
	}
	defer db.Close()
	if err := db.Save(item).Error; err != nil {
		logrus.Warnf("DB.Save Error Conn:%s Type:%s (%s)", r.dbKey, reflect.TypeOf(item).String(), err.Error())
		return cerror.NewCodeError(cerror.DB_ERROR, err)
	}
	return nil
}

func (r *DatabaseRepo) SaveMany(items ...interface{}) error {

	return r.InvokeTransation(func(db *gorm.DB) error {
		tx := db.Begin()
		for _, item := range items {
			if err := tx.Save(item).Error; err != nil {
				logrus.Warnf("DB.SaveMany Error Conn:%s Type:%s (%s)", r.dbKey, reflect.TypeOf(item).String(), err.Error())
				tx.Rollback()
				return err
			}
		}
		tx.Commit()
		return nil
	})
}

func (r *DatabaseRepo) Update(model interface{}, query string, params []interface{}, item ...interface{}) error {

	r.put()
	defer r.pop()

	db, err := GetDB(r.dbKey)
	if err != nil {
		return cerror.NewCodeError(cerror.DB_ERROR, err)

	}
	defer db.Close()
	db = db.Model(model).Where(query, params...)
	if err := db.Update(item...).Error; err != nil {
		logrus.Warnf("DB.Update Error Conn:%s Type:%s (%s)", r.dbKey, reflect.TypeOf(item).String(), err.Error())
		return cerror.NewCodeError(cerror.DB_ERROR, err)

	}
	return nil
}

func (r *DatabaseRepo) UpdateColumn(model interface{}, query string, params []interface{}, attrs ...interface{}) error {

	r.put()
	defer r.pop()

	db, err := GetDB(r.dbKey)
	if err != nil {
		return cerror.NewCodeError(cerror.DB_ERROR, err)

	}
	defer db.Close()
	if err := db.Model(model).Where(query, params...).UpdateColumn(attrs...).Error; err != nil {
		logrus.Warnf("DB.UpdateColumn Error Conn:%s Type:%s (%s)", r.dbKey, reflect.TypeOf(model).String(), err.Error())
		return cerror.NewCodeError(cerror.DB_ERROR, err)
	}
	return nil
}

func (r *DatabaseRepo) Updates(model interface{}, query string, where []interface{}, item map[string]interface{}) error {

	r.put()
	defer r.pop()

	db, err := GetDB(r.dbKey)
	if err != nil {
		return cerror.NewCodeError(cerror.DB_ERROR, err)

	}
	defer db.Close()
	if err := db.Model(model).Where(query, where...).Updates(item).Error; err != nil {
		logrus.Warnf("DB.Updates Error Conn:%s Type:%s (%s)", r.dbKey, reflect.TypeOf(item).String(), err.Error())
		return cerror.NewCodeError(cerror.DB_ERROR, err)

	}
	return nil
}

func (r *DatabaseRepo) Delete(item interface{}) error {

	r.put()
	defer r.pop()

	db, err := GetDB(r.dbKey)
	if err != nil {
		return cerror.NewCodeError(cerror.DB_ERROR, err)

	}
	defer db.Close()
	if err := db.Delete(item).Error; err != nil {
		logrus.Warnf("DB.Delete Error Conn:%s Type:%s (%s)", r.dbKey, reflect.TypeOf(item).String(), err.Error())
		return cerror.NewCodeError(cerror.DB_ERROR, err)
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

func (r *DatabaseRepo) FirstEX(item interface{}, option SearchOption) *gorm.DB {
	r.put()
	defer r.pop()
	db, err := GetDB(r.dbKey)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if option.Order != "" {
		db = db.Order(option.Order)
	}
	values := []interface{}{}
	if option.Where != "" {
		values = append(values, option.Where)
		values = append(values, option.Params...)
	}
	db = db.Offset(option.Offset)
	db = db.First(item, values...)
	if db.Error != nil && !db.RecordNotFound() {
		logrus.Warnf("DB.FirstEX Error Conn:%s Type:%s WHERE:%v (%s)", r.dbKey, reflect.TypeOf(item).String(), option.Where, db.Error.Error())
	}
	return db

}

func (r *DatabaseRepo) FindEx(item interface{}, option SearchOption) *gorm.DB {
	r.put()
	defer r.pop()
	db, err := GetDB(r.dbKey)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if option.Order != "" {
		db = db.Order(option.Order)
	}
	values := []interface{}{}
	if option.Where != "" {
		values = append(values, option.Where)
		values = append(values, option.Params...)
	}

	dbCount := db

	if option.TotalOut != nil {
		if err := dbCount.Model(item).Where(option.Where, option.Params...).Count(option.TotalOut).Error; err != nil {
			logrus.Warnf("DB.FirstEX:Count Error Conn:%s Type:%s WHERE:%v (%s)", r.dbKey, reflect.TypeOf(item).String(), option.Where, err.Error())

		}
	}

	if option.Order != "" {
		db = db.Order(option.Order)
	}

	db = db.Offset(option.Offset)
	if option.Limit > 0 || option.Limit == -1 {
		db = db.Limit(option.Limit)
	}

	db = db.Find(item, values...)
	if db.Error != nil && !db.RecordNotFound() {
		logrus.Warnf("DB.FirstEX Error Conn:%s Type:%s WHERE:%v (%s)", r.dbKey, reflect.TypeOf(item).String(), option.Where, db.Error.Error())
	}
	return db

}

func (r *DatabaseRepo) First(item interface{}, where ...interface{}) *gorm.DB {

	r.put()
	defer r.pop()

	db, err := GetDB(r.dbKey)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db = db.First(item, where...)
	if db.Error != nil && !db.RecordNotFound() {
		logrus.Warnf("DB.First Error Conn:%s Type:%s WHERE:%v (%s)", r.dbKey, reflect.TypeOf(item).String(), where, db.Error.Error())
	}
	return db
}

func (r *DatabaseRepo) DeleteWhere(item interface{}, query string, params ...interface{}) error {

	r.put()
	defer r.pop()

	db, err := GetDB(r.dbKey)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := db.Where(query, params...).Delete(item).Error; err != nil {
		logrus.Warnf("DB.DeleteWhere Error Conn:%s Type:%s WHERE:%s (%s)", r.dbKey, reflect.TypeOf(item).String(), query, err.Error())
		return cerror.NewCodeError(cerror.DB_ERROR, err)
	}
	return nil
}

func (r *DatabaseRepo) Find(offset int, limit int, list interface{}, where ...interface{}) error {

	r.put()
	defer r.pop()

	db, err := GetDB(r.dbKey)
	if err != nil {
	}
	defer db.Close()
	db = db.Offset(offset)
	if limit > 0 {
		db = db.Limit(limit)
	}
	if err := db.Find(list, where...).Error; err != nil && !db.RecordNotFound() {
		logrus.Warnf("DB.Find Error Conn:%s Type:%s WHERE:%v (%s)", r.dbKey, reflect.TypeOf(list).String(), where, err.Error())
		return cerror.NewCodeError(cerror.DB_ERROR, err)
	}
	return nil
}

func (r *DatabaseRepo) FindC(offset int, limit int, list interface{}, where ...interface{}) (int, error) {

	r.put()
	defer r.pop()

	db, err := GetDB(r.dbKey)
	if err != nil {
	}
	defer db.Close()

	db = db.Offset(offset)
	if limit > 0 {
		db = db.Limit(limit)
	}
	if err := db.Find(list, where...).Error; err != nil {
		logrus.Warnf("DB.Find Error Conn:%s Type:%s WHERE:%v (%s)", r.dbKey, reflect.TypeOf(list).String(), where, err.Error())
		return 0, cerror.NewCodeError(cerror.DB_ERROR, err)
	}
	return 0, nil
}

func (r *DatabaseRepo) Count(item interface{}, query string, values ...interface{}) (int, error) {

	r.put()
	defer r.pop()

	db, err := GetDB(r.dbKey)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	total := 0
	if err := db.Model(item).Where(query, values...).Count(&total).Error; err != nil {
		logrus.Warnf("DB.Count Error Conn:%s SQL:%s WHERE:%s... (%s)", r.dbKey, query, values, err.Error())
		return -1, cerror.NewCodeError(cerror.DB_ERROR, err)
	}
	return total, nil
}

func (r *DatabaseRepo) Raw(sql string, values ...interface{}) error {

	r.put()
	defer r.pop()

	db, err := GetDB(r.dbKey)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := db.Raw(sql, values...).Error; err != nil {
		logrus.Warnf("DB.Raw Error Conn:%s SQL:%s WHERE:%s... (%s)", r.dbKey, sql, values, err.Error())
		return cerror.NewCodeError(cerror.DB_ERROR, err)
	}
	return nil
}

type RowScanHandler func(*sql.Rows) error
type TrasnsationInvokeHandler func(db *gorm.DB) error

func (r *DatabaseRepo) InvokeTransation(callback TrasnsationInvokeHandler) error {
	r.put()
	defer r.pop()
	db, err := GetDB(r.dbKey)

	if err != nil {
		panic(err)
	}
	defer db.Close()
	return callback(db)
}

func (r *DatabaseRepo) RawSelect(rawSQL string, rowScanCallback RowScanHandler, values ...interface{}) error {

	r.put()
	defer r.pop()

	db, err := GetDB(r.dbKey)

	if err != nil {
		panic(err)
	}

	defer db.Close()
	var rows *sql.Rows
	if rows, err = db.Raw(rawSQL, values...).Rows(); err != nil {
		logrus.Warnf("DB.RawSelect Error Conn:%s SQL:%s WHERE:%s... (%s)", r.dbKey, rawSQL, values, err.Error())
		return cerror.NewCodeError(cerror.DB_ERROR, err)
	}

	defer rows.Close()
	for rows.Next() {
		if rowScanCallback != nil {
			err = rowScanCallback(rows)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *DatabaseRepo) RawSelectScan(rawSQL string, callback func(*gorm.DB, *sql.Rows) error, values ...interface{}) error {

	r.put()
	defer r.pop()

	db, err := GetDB(r.dbKey)

	if err != nil {
		panic(err)
	}

	defer db.Close()
	var rows *sql.Rows
	if rows, err = db.Raw(rawSQL, values...).Rows(); err != nil {
		logrus.Warnf("DB.RawSelectScan Error Conn:%s SQL:%s WHERE:%s... (%s)", r.dbKey, rawSQL, values, err.Error())
		return cerror.NewCodeError(cerror.DB_ERROR, err)
	}

	defer rows.Close()
	for rows.Next() {
		if e := callback(db, rows); e != nil {
			return e
		}
	}
	return nil
}

func (r *DatabaseRepo) Exec(sql string, values ...interface{}) error {

	r.put()
	defer r.pop()
	db, err := GetDB(r.dbKey)
	if err != nil {
		return err
	}
	defer db.Close()
	db = db.Exec(sql, values...)
	return db.Error
}

func (r *DatabaseRepo) ExecuteScalar(rawSQL string, params []interface{}, values ...interface{}) error {

	r.put()
	defer r.pop()

	db, err2 := GetDB(r.dbKey)

	if err2 != nil {
		return err2
	}

	defer db.Close()
	var rows *sql.Rows
	var err error
	if rows, err = db.Raw(rawSQL, params...).Rows(); err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		rows.Scan(values...)
		break
	}
	return nil
}
