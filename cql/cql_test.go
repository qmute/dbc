package cql

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gocql/gocql"
)

var (
	hosts    = []string{}
	username = ""
	password = ""
	keyspace = "dbc_test"
)

func init() {
	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	hosts = strings.Split(host, ",")
}

func getOption() []Option {
	opts := []Option{
		WithTimeout(500 * time.Millisecond),
		WithProtoVersion(4),
	}

	if username != "" && password != "" {
		opts = append(opts, WithAliyunAuth(username, password))
	}

	return opts
}

func TestConnect(t *testing.T) {
	opts := getOption()
	se, err := Connect(hosts, keyspace, opts...)
	if err != nil {
		t.Errorf("connect error: %v", err)
		return
	}
	defer se.Close()

	t.Logf("session closed: %t", se.Closed())
}

func createKeyspace(se *gocql.Session) error {
	return ExecStmt(se, fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s 
	WITH REPLICATION = { 
		'class' : 'SimpleStrategy',
		'replication_factor' : 1
	}`, keyspace))
}

func TestConnectWithMigration(t *testing.T) {
	mig := &Migration{
		CreateKeyspace: createKeyspace,
		CreateTables:   nil,
	}
	opts := getOption()
	se, err := ConnectWithMigration(hosts, keyspace, mig, opts...)
	if err != nil {
		t.Errorf("connect error: %v", err)
		return
	}
	defer se.Close()

	t.Logf("session closed: %t", se.Closed())

}
