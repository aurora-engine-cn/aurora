package strutils

import (
	"gitee.com/aurora-engine/aurora/utils/uerr"
	"strconv"
)

func IsEmpty(value string) bool {
	return value == ""
}

// Int 数字字符串转化 int
func Int(v string) (Int int) {
	parseInt, err := strconv.ParseInt(v, 0, 0)
	uerr.UtilError(err)
	return int(parseInt)
}

// Float 浮点字符串转化为 float64
func Float(v string) (Float float64) {
	Float, err := strconv.ParseFloat(v, 64)
	uerr.UtilError(err)
	return
}

// String  基础数据类型转化为字符串
func String(v any) (Str string) {
	if v == nil {
		return
	}
	switch v.(type) {
	case string:
		Str = v.(string)
	case int:
		Str = strconv.Itoa(v.(int))
	case float64:
		Str = strconv.FormatFloat(v.(float64), 'f', 2, 64)
	case bool:
		Str = strconv.FormatBool(v.(bool))
	}
	return
}
