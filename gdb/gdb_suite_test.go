package gdb_test

import (
	"context"
	"io"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/qmute/dbc/gdb"
)

func TestGdb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gdb Suite")
}

var cleaner func()
var db *gorm.DB
var ctx context.Context

var _ = BeforeSuite(func() {
	logrus.SetOutput(io.Discard) // 测试期间保持安静

	var err error
	con, f := gdb.DockerPg()
	cleaner = f
	db, err = gdb.ConnectToPG(con.GormConn(), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	Ω(err).To(Succeed())

})

var _ = AfterSuite(func() {
	cleaner()
})

var _ = BeforeEach(func() {
	ctx = context.Background()
})
var _ = AfterEach(func() {
	err := db.Exec(`DROP SCHEMA public CASCADE;CREATE SCHEMA public;`).Error
	Ω(err).To(Succeed())
})
