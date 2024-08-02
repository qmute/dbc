package gdb

import (
	"fmt"
	"os"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DockerPgOpt struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// GormConn to gorm conn string
func (p *DockerPgOpt) GormConn() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		p.Host,
		p.Port,
		p.User,
		p.Database,
		p.Password,
	)
}

func DockerPgV3(img, version string) (*DockerPgOpt, func()) {
	return innerDockerPg(img, version)
}

// DockerPg 测试辅助函数。 利用 dockertest 生成一次性pg实例。
// version , 可选参数，用于指定pg版本， 默认为 "14"
// 返回 gorm 连接对象以及用于清理此实例的 cleanup 函数
func DockerPg(version ...string) (*DockerPgOpt, func()) {
	return innerDockerPg("postgres", version...)
}

func innerDockerPg(img string, version ...string) (*DockerPgOpt, func()) {
	pool, err := dockertest.NewPool("")
	chk(err)

	ver := "14" // default version
	if len(version) > 0 {
		ver = version[0]
	}

	const testDbName = "test"
	const testPasswd = "test"

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: img,
			Tag:        ver,
			Env: []string{
				"POSTGRES_PASSWORD=" + testPasswd,
				"POSTGRES_DB=" + testDbName,
			},
		}, func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.NeverRestart()
		})

	chk(err)

	port := resource.GetPort("5432/tcp")
	host := "localhost"
	if s := os.Getenv("DOCKER_PG_HOST"); s != "" {
		host = s
	}

	opt := &DockerPgOpt{
		Host:     host,
		Port:     port,
		User:     "postgres",
		Password: testPasswd,
		Database: testDbName,
	}

	err = pool.Retry(func() error {
		cfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}
		gdb, err := ConnectToPG(opt.GormConn(), cfg)
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

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
