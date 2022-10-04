package web

import (
	"net/http"
	"reflect"
)

const (
	AuroraMultipartFile = "AuroraMultipartFile" //文件解析
	AuroraRequest       = "AuroraRequest"       //go 原生请求
	AuroraResponse      = "AuroraResponse"      //go 原生响应
	AuroraValues        = "AuroraValues"
)

// Context 上下文参数，主要用于在业务之间传递 数据使用
// 上下文参数中获取请求参数需要依赖于传递的参数名称
// Ctx 不是线程安全的，在请求中出现多线程操作需要使用锁来保证安全性
type Context map[string]interface{}

func (ctx Context) Clear() {
	for key, _ := range ctx {
		delete(ctx, key)
	}
}

// Request 返回元素 Request
func (ctx Context) Request() *http.Request {
	return ctx[AuroraRequest].(*http.Request)
}

// Response 返回元素 ResponseWriter
func (ctx Context) Response() http.ResponseWriter {
	return ctx[AuroraResponse].(http.ResponseWriter)
}

// MultipartFile 获取POST文件上传
func (ctx Context) MultipartFile() *MultipartFile {
	if file, b := ctx[AuroraMultipartFile]; b {
		return file.(*MultipartFile)
	}
	return nil
}

// Return 设置中断处理，多次调用会覆盖之前设置的值
func (ctx Context) Return(value ...interface{}) {
	values := make([]reflect.Value, 0)
	for _, v := range value {
		values = append(values, reflect.ValueOf(v))
	}
	ctx[AuroraValues] = values
}
