package gdb

import (
	"database/sql/driver"

	"github.com/jackc/pgtype"
	"github.com/jinzhu/copier"
)

// IntArray 自定义类型以便支持 []int
// 底层功能代理给 pgtype.Int8Array
type IntArray []int

func (p *IntArray) Scan(src any) error {
	delegate := &pgtype.Int8Array{}
	if err := delegate.Scan(src); err != nil {
		return err
	}
	var dest []int
	if err := delegate.AssignTo(&dest); err != nil {
		return err
	}
	*p = dest
	return nil
}

func (p IntArray) Value() (driver.Value, error) {
	delegate := pgtype.Int8Array{}
	if err := delegate.Set([]int(p)); err != nil {
		return nil, err
	}
	return delegate.Value()
}

// StringArray 自定义类型以便支持 []string
// 底层功能代理给 pgtype.TextArray
type StringArray []string

func (p *StringArray) Scan(src any) error {
	delegate := &pgtype.TextArray{}
	if err := delegate.Scan(src); err != nil {
		return err
	}
	var dest []string
	if err := delegate.AssignTo(&dest); err != nil {
		return err
	}
	*p = dest
	return nil
}

func (p StringArray) Value() (driver.Value, error) {
	delegate := pgtype.TextArray{}
	if err := delegate.Set([]string(p)); err != nil {
		return nil, err
	}
	return delegate.Value()
}

// Float64Array 自定义类型以便支持 []float64
// 底层功能代理给 pgtype.Float8Array
type Float64Array []float64

func (p *Float64Array) Scan(src any) error {
	delegate := &pgtype.Float8Array{}
	if err := delegate.Scan(src); err != nil {
		return err
	}
	var dest []float64
	if err := delegate.AssignTo(&dest); err != nil {
		return err
	}
	*p = dest
	return nil
}

func (p Float64Array) Value() (driver.Value, error) {
	delegate := pgtype.Float8Array{}
	if err := delegate.Set([]float64(p)); err != nil {
		return nil, err
	}
	return delegate.Value()
}

// Copy 用于 model/domain 之间的转换。 可以避免 nil slice 作为 null 保存入库
func Copy(to, from any, converter ...copier.TypeConverter) error {
	opt := copyOption
	if len(converter) > 0 {
		// 支持自定义转换器
		opt.Converters = append(opt.Converters, converter...)
	}
	return copier.CopyWithOption(to, from, opt)
}

// copyOption 配合copier.CopyWithOption
var copyOption = copier.Option{
	Converters: []copier.TypeConverter{
		{
			SrcType: []int{},
			DstType: IntArray{},
			Fn: func(src interface{}) (interface{}, error) {
				s := src.([]int)
				if s == nil {
					s = []int{} // 如果是nil，生成一个空的slice，这样数据库中永远就不会存null
				}
				return IntArray(s), nil
			},
		},
		{
			SrcType: []string{},
			DstType: StringArray{},
			Fn: func(src interface{}) (interface{}, error) {
				s := src.([]string)
				if s == nil {
					s = []string{} // 如果是nil，生成一个空的slice，这样数据库中永远就不会存null
				}
				return StringArray(s), nil
			},
		},
		{
			SrcType: []float64{},
			DstType: Float64Array{},
			Fn: func(src interface{}) (interface{}, error) {
				s := src.([]float64)
				if s == nil {
					s = []float64{} // 如果是nil，生成一个空的slice，这样数据库中永远就不会存null
				}
				return Float64Array(s), nil
			},
		},
	},
}
