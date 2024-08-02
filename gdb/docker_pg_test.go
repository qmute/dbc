package gdb_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"github.com/qmute/dbc/gdb"
)

var _ = Describe("Docker Pg", func() {
	var cleanup func()
	var db *gorm.DB
	BeforeEach(func() {
		con, f := gdb.DockerPg()
		d, err := gdb.ConnectToPG(con.GormConn(), &gorm.Config{})
		Ω(err).To(Succeed())

		cleanup = f
		db = d
	})

	AfterEach(func() { cleanup() })

	DescribeTable("从数据库读取常数",
		func(n int) {
			var tmp struct {
				Num int
			}
			err := db.Raw(fmt.Sprintf(`select %d as num`, n)).Scan(&tmp).Error
			Ω(err).To(Succeed())
			Ω(tmp.Num).To(Equal(n))
		},
		Entry("选出常数5", 5),
		Entry("选出常数8", 8),
	)
})
