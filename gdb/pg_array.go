package gdb

import "github.com/jackc/pgtype"

// Deprecated: 请使用 gdb.IntArray 类型，无需手工转换。此方法将在未来版本中删除
// NewPgInt8Array 创建pg底层int数据
func NewPgInt8Array(src []int) pgtype.Int8Array {
	if src == nil {
		src = []int{}
	}
	result := pgtype.Int8Array{}
	if err := result.Set(src); err != nil {
		panic("pg driver err: " + err.Error())
	}
	return result

}

// PgInt8ArrayToSlice 请使用 gdb.IntArray 类型，无需手工转换。此方法将在未来版本中删除
func PgInt8ArrayToSlice(a pgtype.Int8Array) []int {
	if a.Status == pgtype.Undefined {
		return nil
	}
	var l []int
	if err := a.AssignTo(&l); err != nil {
		panic("pg driver err: " + err.Error())
	}
	return l
}

// Deprecated: 请使用 gdb.StringArray 类型，无需手工转换。此方法将在未来版本中删除
// NewPgTextArray
func NewPgTextArray(src []string) pgtype.TextArray {
	if src == nil {
		src = []string{}
	}
	result := pgtype.TextArray{}
	if err := result.Set(src); err != nil {
		panic("pg driver err: " + err.Error())
	}
	return result
}

// Deprecated: 请使用 gdb.StringArray 类型，无需手工转换。此方法将在未来版本中删除
// PgTextArrayToSlice
func PgTextArrayToSlice(a pgtype.TextArray) []string {
	if a.Status == pgtype.Undefined {
		// 如果里面没值，则直接返回空 todo 看pgx 文档，再确认一次
		return nil
	}
	var l []string
	if err := a.AssignTo(&l); err != nil {
		panic("pg driver err: " + err.Error())
	}
	return l
}

// Deprecated: 使用 serializer:json tag，无需手工转换。此方法将在未来版本中删除
// NewPgJSONB
func NewPgJSONB(src interface{}) pgtype.JSONB {
	result := pgtype.JSONB{}
	if err := result.Set(src); err != nil {
		panic("pg driver err: " + err.Error())
	}

	return result
}

// Deprecated: 使用 serializer:json tag，无需手工转换。此方法将在未来版本中删除
// PgJSONBToInterface
func PgJSONBToInterface(from pgtype.JSONB, to interface{}) {
	if from.Status == pgtype.Undefined || from.Status == pgtype.Null {
		return
	}

	if err := from.AssignTo(to); err != nil {
		panic("pg driver err: " + err.Error())
	}

	return

}
