package gdb

import (
	"database/sql"
	"errors"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect connect to db
func Connect(dialector gorm.Dialector, cfg *gorm.Config, opts ...Option) (*gorm.DB, error) {
	return connect(dialector, cfg, opts...)
}

// ConnectToMysql connect to mysql
func ConnectToMysql(conn string, cfg *gorm.Config, opts ...Option) (*gorm.DB, error) {
	return connect(mysql.Open(conn), cfg, opts...)
}

// ConnectToPG connect to PostgreSQL
func ConnectToPG(conn string, cfg *gorm.Config, opts ...Option) (*gorm.DB, error) {
	return connect(postgres.Open(conn), cfg, opts...)
}

func connect(dialector gorm.Dialector, cfg *gorm.Config, opts ...Option) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, cfg)
	if err != nil {
		return nil, err
	}

	for _, o := range opts {
		o.apply(db)
	}

	return db, nil
}

// NotFound 是否没有找到
func NotFound(err error) bool {
	if err == nil {
		return false
	}

	l := []error{gorm.ErrRecordNotFound, sql.ErrNoRows}

	for _, v := range l {
		if errors.Is(err, v) {
			return true
		}
	}

	// 有时error会被rpc远程传递，变成rpc error，这时只能用字符串判断了
	strList := []string{gorm.ErrRecordNotFound.Error(), sql.ErrNoRows.Error()}
	for _, v := range strList {
		if err.Error() == v {
			return true
		}
	}

	return false
}

// Dup 是否重复
func Dup(err error) bool {
	if err == nil {
		return false
	}

	// todo 使用 mysql.MySQLError 来作精确判断?
	errStr := []string{"Error 1062", "Duplicate entry", "duplicate key value"}
	for _, s := range errStr {
		if strings.Contains(err.Error(), s) {
			return true
		}
	}
	return false
}

// WithTx 事务
func WithTx(db *gorm.DB, f func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return db.Transaction(f, opts...)
}

// Exist 是否存在
func Exist(db *gorm.DB) (bool, error) {
	var n int
	err := db.Select(`1`).Limit(1).Row().Scan(&n)
	if NotFound(err) {
		return false, nil
	}
	return true, err
}
