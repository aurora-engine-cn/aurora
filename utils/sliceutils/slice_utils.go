package sliceutils

import (
	"strconv"
)

// IsEmpty 判断切片是否为空或是否存在元素
// slice 为空 或者 元素为 0 则返回 true
func IsEmpty[T any](slice []T) bool {
	return slice == nil || len(slice) == 0
}

// StrSlice 任意数据转化为字符串切片
// v 为空或 nil 则返回一个空元素的切片
func StrSlice(v ...any) (strings []string) {
	if len(v) == 0 || v == nil {
		return
	}
	for i := 0; i < len(v); i++ {
		value := v[i]
		switch value.(type) {
		case string:
			strings = append(strings, value.(string))
		case int:
			strings = append(strings, strconv.Itoa(value.(int)))
		case float64:
			Float := strconv.FormatFloat(value.(float64), 'f', 2, 64)
			strings = append(strings, Float)
		case bool:
			Bool := strconv.FormatBool(value.(bool))
			strings = append(strings, Bool)
		}
	}
	return
}

// Slice 返回一个任意切片
func Slice[T any](v ...T) (values []T) {
	return v
}
