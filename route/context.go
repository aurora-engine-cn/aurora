package route

import (
	"net/http"
	"reflect"
)

const (
	auroraQueryCache         = "auroraQueryCache"
	auroraFormCache          = "auroraFormCache"
	auroraMaxMultipartMemory = "auroraMaxMultipartMemory"
)

// Middleware 中间件类型
type Middleware func(Ctx) bool

// Ctx 上下文参数，主要用于在业务之间传递 数据使用
// 上下文参数中获取请求参数需要依赖于传递的参数名称
// Ctx 不是线程安全的，在请求中出现多线程操作需要使用锁来保证安全性
type Ctx map[string]interface{}

func (c Ctx) Clear() {
	for key, _ := range c {
		delete(c, key)
	}
}

// Request 返回元素 Request
func (c Ctx) Request() *http.Request {
	return c[request].(*http.Request)
}

// Response 返回元素 ResponseWriter
func (c Ctx) Response() http.ResponseWriter {
	return c[response].(http.ResponseWriter)
}
// Return 设置中断处理，多次调用会覆盖之前设置的值
func (c Ctx) Return(value ...interface{}) {
	values := make([]reflect.Value, 0)
	for _, v := range value {
		values = append(values, reflect.ValueOf(v))
	}
	c["AuroraValues"] = values
}
