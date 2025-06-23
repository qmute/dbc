package gdb_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/qmute/dbc/gdb"
)

var _ = Describe("Types", func() {
	type Tmp struct {
		Id  int              // 系统主键
		Num gdb.IntArray     `gorm:"type:bigint[]"`
		Str gdb.StringArray  `gorm:"type:text[]"`
		F64 gdb.Float64Array `gorm:"type:float[]"`
	}

	type DTmp struct {
		Id  int // 系统主键
		Num []int
		Str []string
		F64 []float64
	}

	BeforeEach(func() {
		err := db.AutoMigrate(&Tmp{})
		Ω(err).To(Succeed())

		err = db.Create(&Tmp{
			Num: gdb.IntArray{1, 2, 3},
			Str: gdb.StringArray{"a", "b", "c"},
			F64: gdb.Float64Array{1.1, 2.2, 3.3, 8.99},
		}).Error

		Ω(err).To(Succeed())
	})

	It("Scan", func() {
		tmp := &Tmp{}
		err := db.Take(tmp).Error
		Ω(err).To(Succeed())

		Ω(tmp.Num).To(Equal(gdb.IntArray{1, 2, 3}))
		Ω(tmp.Str).To(Equal(gdb.StringArray{"a", "b", "c"}))
		Ω(tmp.F64).To(Equal(gdb.Float64Array{1.1, 2.2, 3.3, 8.99}))

		dtmp := &DTmp{}
		err = gdb.Copy(dtmp, tmp)
		Ω(err).To(Succeed())
		Ω(dtmp).To(BeEquivalentTo(&DTmp{
			Id:  tmp.Id,
			Num: tmp.Num,
			Str: tmp.Str,
			F64: tmp.F64,
		}))

		Ω(dtmp.Num).To(Equal([]int{1, 2, 3}))
		Ω(dtmp.Str).To(Equal([]string{"a", "b", "c"}))
		Ω(dtmp.F64).To(Equal([]float64{1.1, 2.2, 3.3, 8.99}))

	})
})
