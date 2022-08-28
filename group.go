package aurora

import (
	"net/http"
	"strings"
)

/*
group 路由分组
初始化的 分组变量不会携带全局的Use
group 可以设定局部的全局Use
*/
type Group struct {
	a          *Engine
	middleware []Middleware
	prefix     string
}

// Get 请求
func (g *Group) Get(url string, control Controller, middleware ...Middleware) {
	middleware = mergeMiddleware(g.middleware, middleware)
	g.a.register(http.MethodGet, g.prefix+url, control, middleware...)
}

// Post 请求
func (g *Group) Post(url string, control Controller, middleware ...Middleware) {
	middleware = mergeMiddleware(g.middleware, middleware)
	g.a.register(http.MethodPost, g.prefix+url, control, middleware...)
}

// Put 请求
func (g *Group) Put(url string, control Controller, middleware ...Middleware) {
	middleware = mergeMiddleware(g.middleware, middleware)
	g.a.register(http.MethodPut, g.prefix+url, control, middleware...)
}

// Delete 请求
func (g *Group) Delete(url string, control Controller, middleware ...Middleware) {
	middleware = mergeMiddleware(g.middleware, middleware)
	g.a.register(http.MethodDelete, g.prefix+url, control, middleware...)
}

// Head 请求
func (g *Group) Head(url string, control Controller, middleware ...Middleware) {
	middleware = mergeMiddleware(g.middleware, middleware)
	g.a.register(http.MethodHead, g.prefix+url, control, middleware...)
}

// Group 路由分组  必须以 “/” 开头分组
func (g *Group) Group(url string, middleware ...Middleware) *Group {
	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}
	//多次分组 叠加之前的分组处理器
	h := mergeMiddleware(g.middleware, middleware)
	return &Group{
		prefix:     g.prefix + url,
		a:          g.a,
		middleware: h,
	}
}

// Use 基于 group 的分组添加 Middleware
func (g *Group) Use(middleware ...Middleware) {
	if g.middleware == nil {
		g.middleware = middleware
		return
	}
	g.middleware = append(g.middleware, middleware...)
}

// mergeMiddleware 合并两个 Middleware
// g 分组全局 Middleware
// h 局部 Middleware
func mergeMiddleware(g, h []Middleware) []Middleware {
	if g == nil && h == nil {
		return nil
	}
	middleware := make([]Middleware, 0)
	middleware = append(middleware, g...)
	middleware = append(middleware, h...)
	return middleware
}
