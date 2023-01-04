package aurora

import (
	"gitee.com/aurora-engine/aurora/web"
	"net/http"
	"strings"
)

// Get 请求
func (engine *Engine) Get(url string, control any, middleware ...web.Middleware) {
	engine.register(http.MethodGet, url, control, middleware...)
}

// Post 请求
func (engine *Engine) Post(url string, control any, middleware ...web.Middleware) {
	engine.register(http.MethodPost, url, control, middleware...)
}

// Put 请求
func (engine *Engine) Put(url string, control any, middleware ...web.Middleware) {
	engine.register(http.MethodPut, url, control, middleware...)
}

// Delete 请求
func (engine *Engine) Delete(url string, control any, middleware ...web.Middleware) {
	engine.register(http.MethodDelete, url, control, middleware...)
}

// Head 请求
func (engine *Engine) Head(url string, control any, middleware ...web.Middleware) {
	engine.register(http.MethodHead, url, control, middleware...)
}

// register 通用注册器
func (engine *Engine) register(method string, url string, control any, middleware ...web.Middleware) {
	engine.Route.Cache(method, url, control, middleware...)
}

// Group 路由分组  必须以 “/” 开头分组
// Group 和 Aurora 都有 相同的 http 方法注册
func (engine *Engine) Group(url string, middleware ...web.Middleware) *Group {
	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}
	//分组处理的 handles 和 群居的 handle 是区分开的，该处handle只作用于通过该分组创建的 接口，在调用接口之前该 handles会被执行
	return &Group{
		prefix:     url,
		a:          engine,
		middleware: middleware,
	}
}
