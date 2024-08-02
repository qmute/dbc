package cql

import "github.com/gocql/gocql"

// MigrateFunc 数据初始化func
type MigrateFunc func(se *gocql.Session) error

// Migration 数据初始化结构体
type Migration struct {
	CreateKeyspace MigrateFunc
	CreateTables   MigrateFunc
}
