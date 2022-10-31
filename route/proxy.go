package route

import (
	"bytes"
	"gitee.com/aurora-engine/aurora/utils/stringutils"
	"gitee.com/aurora-engine/aurora/web"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"reflect"
	"strings"
)

var errImp = reflect.TypeOf(new(error)).Elem()

// Proxy
// 主要完成 接口调用 和 处理响应
type Proxy struct {
	*Router
	errType     reflect.Type
	Rew         http.ResponseWriter // http 响应体
	Req         *http.Request       // http 请求体
	Middleware  []web.Middleware    // 当前控制器待执行中间件
	Context     web.Context         // 请求上下文数据
	File        *web.MultipartFile  // Post文件上传解析值
	Control     Controller          // 待处理执行器
	values      []reflect.Value     // 处理器返回值
	UrlVariable []string            // RESTFul  顺序值
	RESTFul     map[string]any      // RESTFul  K/V值
	view        ViewHandle          // 支持自定义视图渲染机制
	Recover     WebRecover          // 路由错误捕捉器
}

// start 路由查询入口
func (proxy *Proxy) start() {
	//defer 处理在执行接口期间的一切 panic处理
	defer proxy.Recover(proxy)
	//存储error类型 用于catch捕捉
	proxy.errType = errImp
	c := proxy.Control
	c.Proxy = proxy
	c.Context = proxy.Context
	c.UrlVariable = proxy.UrlVariable
	c.RESTFul = proxy.RESTFul
	//执行中间件
	middlewares := proxy.Middleware
	for _, middleware := range middlewares {
		if b := middleware(proxy.Context); !b {
			goto end
		}
	}
	// 请求参数解析
	c.analysisInput(proxy.Req)
	// 执行请求方法
	proxy.values = c.invoke()

end: // end 执行 结果处理，主要方便 在中间件中发生中断，设置响应策略
	// 判断防止中间件误用中断 并且放行 如果处理器执行了 中断将失效
	if proxy.values == nil {
		//如果中间件中断，尝试从中间件中拿到一个结果
		if v, b := proxy.Context["AuroraValues"]; b {
			proxy.values = v.([]reflect.Value)
		}
	}
	proxy.resultHandler()
	return
}

// resultHandler 接口返回值相应处理
// 更具返回的 数据类型对应如何返回数据给前端
func (proxy *Proxy) resultHandler() {
	if proxy.values == nil {
		return
	}
	header := proxy.Rew.Header()
	// 处理响应结果之前 判空 Content-Type，以支持用户自定义返回格式,默认会添加json格式
	get := header.Get(contentType)
	if stringutils.IsEmpty(get) {
		header.Set(contentType, ResourceMapType[".json"])
	}
	for i := 0; i < len(proxy.values); i++ {
		switch proxy.values[i].Kind() {
		case reflect.String:
			value := proxy.values[i].Interface().(string)
			stringData(proxy, value)
		//对于接口的返回，目前只做了对错误的支持，app 开发中对抽象类型的设计应该不会太多，大部分直接返回实体数据了
		case reflect.Pointer, reflect.Struct, reflect.Slice, reflect.Int, reflect.Float64, reflect.Bool, reflect.Map:
			otherData(proxy, proxy.values[i])
		case reflect.Interface:
			anyData(proxy, proxy.values[i])
		}
	}
}

// catchError 处理错误捕捉
func (proxy *Proxy) catchError(errType reflect.Type, errValue reflect.Value) {
	v := errValue.Interface()
	if catch, b := proxy.Catches[errType]; b {
		//如果进行了错误捕捉,由于这个请求在同一个处理内，可以选择覆盖掉之前的返回内容，然后继续通过处理方法对返回值进行处理, 在错误处理器中通常不应该再次返回错误
		proxy.values = catch.invoke(errValue)
		proxy.resultHandler()
		return
	} else {
		//没有注册捕捉器的错误均已默认输出控制台的方式展示, 非 panic 方式的错误会被直接输出控制台 error 消息
		switch v.(type) {
		case error:
			proxy.Error(v.(error).Error())
		}
	}
}

// 接口返回string类型处理函数
func stringData(proxy *Proxy, value string) {
	if strings.HasSuffix(value, ".html") {
		HtmlPath := proxy.PathPool.Get().(*bytes.Buffer)
		if value[:1] == "/" {
			value = value[1:]
		}
		//拼接项目路径
		HtmlPath.WriteString(proxy.Root)
		//拼接 静态资源路径 默认情况下为 '/'
		HtmlPath.WriteString(proxy.Resource)
		//拼接 资源真实路径
		HtmlPath.WriteString(value)
		//得到完整 html 页面资源path
		html := HtmlPath.String()
		HtmlPath.Reset() //清空buffer，以便下次利用
		proxy.PathPool.Put(HtmlPath)
		proxy.Rew.Header().Set(contentType, ResourceMapType[".html"])
		proxy.view(html, proxy.Rew, proxy.Context) //视图解析 响应 html 页面
		return
	}
	//处理转发，重定向本质重新走一边路由，找到对应处理的方法
	if strings.HasPrefix(value, "forward:") {
		value = value[8:]
		//请求转发 会携带当前的 请求体 和上下文参数
		c, u, args, ctx := proxy.urlRouter(proxy.Req.Method, value, proxy.Rew, proxy.Req, proxy.Context)
		proxy.handle(c, u, args, proxy.Rew, proxy.Req, ctx)
		return
	}
	proxy.Rew.Write([]byte(value))
}

func otherData(proxy *Proxy, value reflect.Value) {
	of := value.Type()
	var v = value.Interface()
	switch v.(type) {
	case error:
		proxy.catchError(of, value)
		return

	case int, float64, bool:

	default:
		marshal, err := jsoniter.Marshal(value.Interface())
		ErrorMsg(err)
		proxy.Rew.Write(marshal)
	}
}

func anyData(proxy *Proxy, value reflect.Value) {
	valuer := value.Elem()
	of := value.Type()
	var marshal []byte
	var v = valuer.Interface()
	switch v.(type) {
	case error:
		proxy.catchError(of, value)
	case string:
		//对字符串不仅处理
		marshal = []byte(v.(string))
	default:
		s, err := jsoniter.Marshal(v)
		ErrorMsg(err)
		marshal = s
		proxy.Rew.Write(marshal)
		return
	}
}
