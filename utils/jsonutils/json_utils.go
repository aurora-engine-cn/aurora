package jsonutils

import jsoniter "github.com/json-iterator/go"

// JsonString 将数据转化为json字符串
func JsonString(data any) string {
	marshal, err := jsoniter.Marshal(data)
	if err != nil {
		panic(err)
	}
	return string(marshal)
}
