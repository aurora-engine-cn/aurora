package jsonutils

import (
	"gitee.com/aurora-engine/aurora/utils/uerr"
	jsoniter "github.com/json-iterator/go"
)

// Json 将数据转化为json字符串
func Json(data any) string {
	marshal, err := jsoniter.Marshal(data)
	uerr.UtilError(err)
	return string(marshal)
}
