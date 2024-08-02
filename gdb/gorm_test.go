package gdb

import (
	"os"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	conn string
)

func init() {
	conn = os.Getenv("CONN")
}

func TestConnectToMysql(t *testing.T) {
	if true {
		return
	}

	opts := getOptions()
	db, err := ConnectToMysql(conn, &gorm.Config{}, opts...)
	if err != nil {
		t.Fatal(err)
		return
	}

	if err := operation(db); err != nil {
		t.Fatal(err)
		return
	}

	t.Log("connect mysql success")
}

func TestConnectToPG(t *testing.T) {
	if true {
		return
	}

	cfg := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
	}

	opts := getOptions()

	db, err := ConnectToPG(conn, cfg, opts...)
	if err != nil {
		t.Fatal(err)
		return
	}

	if err := operation(db); err != nil {
		t.Fatal(err)
		return
	}

	t.Log("connect pg success")
}

func TestConnectWithPG(t *testing.T) {
	if true {
		return
	}

	opts := getOptions()
	cfg := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
	}

	db, err := Connect(postgres.Open(conn), cfg, opts...)
	if err != nil {
		t.Fatal(err)
		return
	}

	if err := operation(db); err != nil {
		t.Fatal(err)
		return
	}

	t.Log("connect success")

}

func TestConnectWithMysql(t *testing.T) {
	if true {
		return
	}
	opts := getOptions()
	db, err := Connect(mysql.Open(conn), &gorm.Config{}, opts...)
	if err != nil {
		t.Fatal(err)
		return
	}

	if err := operation(db); err != nil {
		t.Fatal(err)
		return
	}

	t.Log("connect success")

}

func getOptions() []Option {
	opts := []Option{
		WithConnMaxLifetime(3 * time.Minute),
		WithMaxIdleConns(10),
		WithMaxOpenConns(50),
		WithPing(10 * time.Second),
	}

	return opts
}

func operation(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	return nil

}
