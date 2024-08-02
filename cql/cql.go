package cql

import (
	"errors"
	"time"

	"github.com/gocql/gocql"
)

var (
	// ErrKeyspaceEmpty 返回一个 "keyspace empty error"
	ErrKeyspaceEmpty = errors.New("keyspace must be not empty")
)

// Connect 返回数据库会话实例.
func Connect(hosts []string, keyspace string, opts ...Option) (*gocql.Session, error) {
	return connect(hosts, keyspace, nil, opts...)
}

// ConnectWithMigration 初始化数据后, 返回数据库会话实例
func ConnectWithMigration(hosts []string, keyspace string, migration *Migration, opts ...Option) (*gocql.Session, error) {
	return connect(hosts, keyspace, migration, opts...)
}

func connect(hosts []string, keyspace string, migration *Migration, opts ...Option) (*gocql.Session, error) {
	if keyspace == "" {
		return nil, ErrKeyspaceEmpty
	}

	cluster := createCluster(hosts, keyspace, opts...)
	return createSessionFromCluster(cluster, keyspace, migration)
}

func createCluster(hosts []string, keyspace string, opts ...Option) *gocql.ClusterConfig {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	for _, o := range opts {
		o.apply(cluster)
	}
	return cluster
}

func createSessionFromCluster(cluster *gocql.ClusterConfig, keyspace string, migration *Migration) (*gocql.Session, error) {
	if err := startMigrate(cluster, migration); err != nil {
		return nil, err
	}

	cluster.Keyspace = keyspace
	return cluster.CreateSession()
}

func createMigrateSession(cluster *gocql.ClusterConfig, keyspace string) (*gocql.Session, error) {
	c := *cluster
	c.Keyspace = keyspace
	c.Timeout = 30 * time.Second
	se, err := c.CreateSession()
	if err != nil {
		return nil, err
	}

	return se, nil
}

func startMigrate(cluster *gocql.ClusterConfig, migration *Migration) error {
	if migration == nil {
		return nil
	}

	if migration.CreateKeyspace != nil {
		se, err := createMigrateSession(cluster, "")
		if err != nil {
			return err
		}
		defer se.Close()

		if err := migration.CreateKeyspace(se); err != nil {
			return err
		}
	}

	if migration.CreateTables != nil {
		se, err := createMigrateSession(cluster, cluster.Keyspace)
		if err != nil {
			return err
		}
		defer se.Close()

		if err := migration.CreateTables(se); err != nil {
			return err
		}
	}

	return nil
}

// ExecStmt 执行一个cql语句.
func ExecStmt(s *gocql.Session, stmt string) error {
	q := s.Query(stmt).RetryPolicy(nil)
	defer q.Release()
	return q.Exec()
}

// NotFound 判断错误是不是 gocql.ErrNotFound
func NotFound(err error) bool {
	return errors.Is(err, gocql.ErrNotFound)
}
