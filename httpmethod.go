package aurora

import (
	"net/http"
	reflect "reflect"
	"strings"
)

type Controller = interface{}

// Get 请求
func (a *Aurora) Get(url string, control Controller, middleware ...Middleware) {
	a.register(http.MethodGet, url, control, middleware...)
}

// Post 请求
func (a *Aurora) Post(url string, control Controller, middleware ...Middleware) {
	a.register(http.MethodPost, url, control, middleware...)
}

// Put 请求
func (a *Aurora) Put(url string, control Controller, middleware ...Middleware) {
	a.register(http.MethodPut, url, control, middleware...)
}

// Delete 请求
func (a *Aurora) Delete(url string, control Controller, middleware ...Middleware) {
	a.register(http.MethodDelete, url, control, middleware...)
}

// Head 请求
func (a *Aurora) Head(url string, control Controller, middleware ...Middleware) {
	a.register(http.MethodHead, url, control, middleware...)
}

// Url 结构体专属 注册器
func (a *Aurora) Url(url string, control Controller, middleware ...Middleware) {
	structUrl := a.analysisStructUrl(url, control, middleware...)
	if a.api == nil {
		a.api = make(map[string][]controlInfo)
	}
	for k, infos := range structUrl {
		if _, b := a.api[k]; !b {
			a.api[k] = infos
		} else {
			a.api[k] = append(a.api[k], infos...)
		}
	}
}

// register 通用注册器
func (a *Aurora) register(method string, url string, control Controller, middleware ...Middleware) {
	if a.api == nil {
		a.api = make(map[string][]controlInfo)
	}
	apis := a.analysisStruct(url, control, middleware...)
	// 查重校验
	//api := controlInfo{path: url, control: control, middleware: middleware}
	if _, b := a.api[method]; !b {
		a.api[method] = make([]controlInfo, 0)
		a.api[method] = append(a.api[method], apis...)
	} else {
		a.api[method] = append(a.api[method], apis...)
	}
}

// Group 路由分组  必须以 “/” 开头分组
// Group 和 Aurora 都有 相同的 http 方法注册
func (a *Aurora) Group(url string, middleware ...Middleware) *Group {
	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}
	//分组处理的 handles 和 群居的 handle 是区分开的，该处handle只作用于通过该分组创建的 接口，在调用接口之前该 handles会被执行
	return &Group{
		prefix:     url,
		a:          a,
		middleware: middleware,
	}
}

// analysisStruct 解析结构体作为接口注册
// control 为结构体或者指向结构体的指针 或者是个函数
func (a *Aurora) analysisStruct(path string, control Controller, middleware ...Middleware) []controlInfo {
	if control == nil {
		return nil
	}
	arr := make([]controlInfo, 0)
	of := reflect.ValueOf(control)
	switch of.Kind() {
	case reflect.Ptr:
		elem := of.Elem()
		if elem.Kind() == reflect.Ptr {
			return nil
		}
		if elem.Kind() != reflect.Struct {
			return nil
		}
		methods := of.NumMethod()
		for i := 0; i < methods; i++ {
			method := of.Type().Method(i)
			//if !method.IsExported() {
			//	continue
			//}
			register := urlRegister(method.Name)
			fun := of.Method(i)
			info := controlInfo{path: path + "/" + register, control: fun.Interface(), middleware: middleware}
			arr = append(arr, info)
		}
		a.control(control)
	case reflect.Struct:
		methods := of.NumMethod()
		for i := 0; i < methods; i++ {
			method := of.Type().Method(i)
			//if !method.IsExported() {
			//	continue
			//}
			register := urlRegister(method.Name)
			fun := of.Method(i)
			info := controlInfo{path: path + "/" + register, control: fun.Interface(), middleware: middleware}
			arr = append(arr, info)
		}
		a.control(control)
	case reflect.Func:
		info := controlInfo{path: path, control: control, middleware: middleware}
		arr = append(arr, info)
	}
	return arr
}

func (a *Aurora) analysisStructUrl(path string, control Controller, middleware ...Middleware) map[string][]controlInfo {
	if control == nil {
		return nil
	}
	httpmethods := make(map[string][]controlInfo)
	of := reflect.ValueOf(control)
	switch of.Kind() {
	case reflect.Ptr:
		elem := of.Elem()
		if elem.Kind() == reflect.Ptr {
			return nil
		}
		if elem.Kind() != reflect.Struct {
			return nil
		}
		methods := of.NumMethod()
		for i := 0; i < methods; i++ {
			method := of.Type().Method(i)
			//if !method.IsExported() {
			//	continue
			//}
			register := urlRegister(method.Name)
			index := strings.Index(register, "/")
			s := register[:index]
			register = register[index+1:]
			m := strings.ToUpper(s)
			if m == http.MethodPost || m == http.MethodGet || m == http.MethodPut || m == http.MethodDelete {
				if _, b := httpmethods[m]; !b {
					httpmethods[m] = make([]controlInfo, 0)
				}
				fun := of.Method(i)
				info := controlInfo{path: path + "/" + register, control: fun.Interface(), middleware: middleware}
				httpmethods[m] = append(httpmethods[m], info)
			}
		}
		a.control(control)
	case reflect.Struct:
		methods := of.NumMethod()
		for i := 0; i < methods; i++ {
			method := of.Type().Method(i)
			//if !method.IsExported() {
			//	continue
			//}
			register := urlRegister(method.Name)
			index := strings.Index(register, "/")
			s := register[:index]
			register = register[index+1:]
			m := strings.ToUpper(s)
			if m == http.MethodPost || m == http.MethodGet || m == http.MethodPut || m == http.MethodDelete {
				if _, b := httpmethods[m]; !b {
					httpmethods[m] = make([]controlInfo, 0)
				}
				fun := of.Method(i)
				info := controlInfo{path: path + "/" + register, control: fun.Interface(), middleware: middleware}
				httpmethods[m] = append(httpmethods[m], info)
			}

		}
		a.control(control)
	}
	return httpmethods
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
