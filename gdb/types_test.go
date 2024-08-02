package gdb_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/qmute/dbc/gdb"
)

var _ = Describe("Types", func() {
	type Tmp struct {
		Id  int             // 系统主键
		Num gdb.IntArray    `gorm:"type:bigint[]"`
		Str gdb.StringArray `gorm:"type:text[]"`
	}

	BeforeEach(func() {
		err := db.AutoMigrate(&Tmp{})
		Ω(err).To(Succeed())

		err = db.Create(&Tmp{
			Num: gdb.IntArray{1, 2, 3},
			Str: gdb.StringArray{"a", "b", "c"},
		}).Error

		Ω(err).To(Succeed())
	})

	It("Scan", func() {
		tmp := &Tmp{}
		err := db.Take(tmp).Error
		Ω(err).To(Succeed())

		Ω(tmp.Num).To(Equal(gdb.IntArray{1, 2, 3}))
		Ω(tmp.Str).To(Equal(gdb.StringArray{"a", "b", "c"}))
	})
})
