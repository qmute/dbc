package gdb

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DockerMysqlOpt struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// GormConn to gorm conn string
func (p DockerMysqlOpt) GormConn() string {
	return fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.Database,
	)
}

// DockerMySQL 测试辅助函数。 利用 dockertest 生成一次性 mysql 实例。
// version , 可选参数，用于指定 mysql 版本， 默认为 "8"
// 返回 gorm 连接对象以及用于清理此实例的 cleanup 函数
func DockerMySQL(version ...string) (*DockerMysqlOpt, func()) {
	return innerDockerMySQL("mysql", version...)
}

func innerDockerMySQL(img string, version ...string) (*DockerMysqlOpt, func()) {
	pool, err := dockertest.NewPool("")
	chk(err)

	ver := "8" // default version
	if len(version) > 0 {
		ver = version[0]
	}

	const (
		testDbName = "test"
		testPasswd = "test"
	)

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: img,
			Tag:        ver,
			Env: []string{
				"MYSQL_ROOT_PASSWORD=" + testPasswd,
				"MYSQL_DATABASE=" + testDbName,
			},
		}, func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.NeverRestart()
		})

	chk(err)

	port := resource.GetPort("3306/tcp")
	host := "localhost"
	if s := os.Getenv("DOCKER_MYSQL_HOST"); s != "" {
		host = s
	}

	user := "root"
	opt := &DockerMysqlOpt{
		Host:     host,
		Port:     port,
		User:     user,
		Password: testPasswd,
		Database: testDbName,
	}

	err = pool.Retry(func() error {
		err := mysql.SetLogger(log.New(io.Discard, "", log.LstdFlags))
		chk(err)
		cfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}
		gdb, err := ConnectToMysql(opt.GormConn(), cfg)
		if err != nil {
			return err
		}

		db, err := gdb.DB()
		if err != nil {
			return err
		}
		return db.Ping()
	})
	chk(err)

	return opt, func() {
		chk(resource.Close())
	}
}
