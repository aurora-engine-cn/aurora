package web

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"runtime/debug"
)

// Recover 用于处理服务器中出现的 panic 消息自定义
type Recover func(w http.ResponseWriter)

// ErrRecover 全局错误 panic 处理
func ErrRecover(w http.ResponseWriter) {
	if v := recover(); v != nil {
		var msg string
		switch v.(type) {
		case error:
			msg = v.(error).Error()
		case string:
			msg = v.(string)
		default:
			marshal, err := jsoniter.Marshal(v)
			if err != nil {
				msg = err.Error()
			}
			msg = string(marshal)
		}
		fmt.Println(msg)
		debug.PrintStack()
		http.Error(w, msg, 500)
		return
	}
}
