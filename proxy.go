package aurora

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
)

const aurora = "aurora"

// Proxy
// 主要完成 接口调用 和 处理响应
type Proxy struct {
	*Engine
	Rew         http.ResponseWriter
	Req         *http.Request
	errType     reflect.Type
	middleware  []Middleware
	Ctx         Ctx
	File        *MultipartFile
	control     controller
	values      []reflect.Value
	UrlVariable []string
	args        map[string]interface{} //REST API 参数解析
	view        views                  //支持自定义视图渲染机制
}

// start 路由查询入口
func (sp *Proxy) start() []reflect.Value {
	//defer 处理在执行接口期间的一切 panic处理
	defer errRecover(sp)
	//存储error类型 用于catch捕捉
	ef := reflect.TypeOf(new(error)).Elem()
	sp.errType = ef
	c := sp.control
	c.p = sp
	c.ctx = sp.Ctx
	c.UrlVariable = sp.UrlVariable
	c.RESTFul = sp.args
	//执行中间件
	middlewares := sp.middleware
	for _, middleware := range middlewares {
		if b := middleware(sp.Ctx); !b {
			goto end
		}
	}
	// 请求参数解析
	c.analysisInput(sp.Req, sp.Rew, sp.Ctx)

	// 执行请求方法
	sp.values = c.invoke()

	// end 执行 结果处理，主要方便 在中间件中发生中断，设置响应策略
end:
	// 判断防止中间件误用中断 并且放行 如果处理器执行了 中断将失效
	if sp.values == nil {
		//如果中间件中断，尝试从中间件中拿到一个结果
		if v, b := sp.Ctx["AuroraValues"]; b {
			sp.values = v.([]reflect.Value)
		}
	}
	sp.resultHandler()
	return nil
}

// resultHandler 接口返回值相应处理
// 更具返回的 数据类型对应如何返回数据给前端
func (sp *Proxy) resultHandler() {
	if sp.values == nil {
		return
	}
	header := sp.Rew.Header()
	if header.Get(contentType) == "" {
		header.Set(contentType, sp.Engine.resourceMapType[".json"])
	}
	for i := 0; i < len(sp.values); i++ {
		v := sp.values[i].Interface()
		//存在 nil 不处理
		if v == nil {
			continue
		}
		switch sp.values[i].Kind() {
		case reflect.String:
			value := sp.values[i].Interface().(string)
			stringData(sp, value)
		//对于接口的返回，目前只做了对错误的支持，web 开发中对抽象类型的设计应该不会太多，大部分直接返回实体数据了
		case reflect.Ptr, reflect.Struct, reflect.Slice, reflect.Int, reflect.Float64, reflect.Bool, reflect.Map:
			otherData(sp, sp.values[i])
		case reflect.Interface:
			anyData(sp, sp.values[i])
		}
	}
}

// catchError 处理错误捕捉
func (sp *Proxy) catchError(errType reflect.Type, errValue reflect.Value) {
	v := errValue.Interface()
	if catch, b := sp.router.catch[errType]; b {
		//如果进行了错误捕捉,由于这个请求在同一个处理内，可以选择覆盖掉之前的返回内容，然后继续通过处理方法对返回值进行处理, 在错误处理器中通常不应该再次返回错误
		sp.values = catch.invoke(errValue)
		sp.resultHandler()
		return
	} else {
		//没有注册捕捉器的错误均已默认输出控制台的方式展示, 非 panic 方式的错误会被直接输出控制台 error 消息
		switch v.(type) {
		case error:
			sp.Error(v.(error).Error())
		}
	}
}

// 接口返回string类型处理函数
func stringData(sp *Proxy, value string) {
	if strings.HasSuffix(value, ".html") {
		HtmlPath := sp.pathPool.Get().(*bytes.Buffer)
		if value[:1] == "/" {
			value = value[1:]
		}
		//拼接项目路径
		HtmlPath.WriteString(sp.projectRoot)
		//拼接 静态资源路径 默认情况下为 '/'
		HtmlPath.WriteString(sp.resource)
		//拼接 资源真实路径
		HtmlPath.WriteString(value)
		//得到完整 html 页面资源path
		html := HtmlPath.String()
		HtmlPath.Reset() //清空buffer，以便下次利用
		sp.pathPool.Put(HtmlPath)
		sp.Rew.Header().Set(contentType, sp.resourceMapType[".html"])
		sp.view.view(html, sp.Rew, nil) //视图解析 响应 html 页面
		return
	}
	//处理转发，重定向本质重新走一边路由，找到对应处理的方法
	if strings.HasPrefix(value, "forward:") {
		value = value[8:]
		//请求转发 会携带当前的 请求体 和上下文参数
		c, u, args, ctx := sp.router.urlRouter(sp.Req.Method, value, sp.Rew, sp.Req, sp.Ctx)
		sp.handle(c, u, args, sp.Rew, sp.Req, ctx)
		return
	}
	sp.Rew.Write([]byte(value))
}

func otherData(sp *Proxy, value reflect.Value) {
	of := value.Type()
	if of.Implements(sp.errType) {
		//错误捕捉
		sp.catchError(of, value)
		return
	}
	marshal, err := json.Marshal(value.Interface())
	ErrorMsg(err)
	sp.Rew.Write(marshal)
}

func anyData(sp *Proxy, value reflect.Value) {
	valuer := value.Elem()
	of := value.Type()
	if !of.Implements(sp.errType) {
		//没有实现,反射校验接口是否实现的小坑，实现接口的形式要和统一，比如 反射类型是指针，实现接口绑定的方式要是指针
		//此处可能返回 interface{} 的数据 没有实现error的当作数据 返回
		var marshal []byte
		v := valuer.Interface()
		switch v.(type) {
		case string:
			//对字符串不仅处理
			marshal = []byte(v.(string))
		default:
			s, err := json.Marshal(v)
			ErrorMsg(err)
			marshal = s
		}
		sp.Rew.Write(marshal)
		return
	}
	//错误捕捉
	sp.catchError(of, value)
}
