package aurora

import (
	"net/http"
	"strings"
)

// Get 请求
func (engine *Engine) Get(url string, control any, middleware ...Middleware) {
	engine.register(http.MethodGet, url, control, middleware...)
}

// Post 请求
func (engine *Engine) Post(url string, control any, middleware ...Middleware) {
	engine.register(http.MethodPost, url, control, middleware...)
}

// Put 请求
func (engine *Engine) Put(url string, control any, middleware ...Middleware) {
	engine.register(http.MethodPut, url, control, middleware...)
}

// Delete 请求
func (engine *Engine) Delete(url string, control any, middleware ...Middleware) {
	engine.register(http.MethodDelete, url, control, middleware...)
}

// Head 请求
func (engine *Engine) Head(url string, control any, middleware ...Middleware) {
	engine.register(http.MethodHead, url, control, middleware...)
}

// register 通用注册器
func (engine *Engine) register(method string, url string, control any, middleware ...Middleware) {
	//if engine.api == nil {
	//	engine.api = make(map[string][]controlInfo)
	//}
	//api := controlInfo{path: url, control: control, middleware: middleware}
	//if _, b := engine.api[method]; !b {
	//	engine.api[method] = make([]controlInfo, 0)
	//	engine.api[method] = append(engine.api[method], api)
	//} else {
	//	engine.api[method] = append(engine.api[method], api)
	//}
	engine.router.Cache(method, url, control, middleware...)
}

// Group 路由分组  必须以 “/” 开头分组
// Group 和 Aurora 都有 相同的 http 方法注册
func (engine *Engine) Group(url string, middleware ...Middleware) *Group {
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

func urlRegister(p string) string {
	if len(p) == 1 {
		return strings.ToLower(p)
	}
	builder := strings.Builder{}
	index := 0
	for i := 1; i < len(p); i++ {
		c := p[i : i+1]
		if c == "_" {
			builder.WriteString(p[index:i] + "/")
			index = i + 1
		}
		if c >= "A" && c <= "Z" {
			builder.WriteString(p[index:i] + "/")
			index = i
		}
	}
	// 处理 最后一个驼峰
	if index != len(p)-1 {
		builder.WriteString(p[index:])
	}
	path := builder.String()
	return strings.ToLower(path)
}
