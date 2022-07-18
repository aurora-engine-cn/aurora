package aurora

import (
	"errors"
	"gitee.com/aurora-engine/aurora/utils"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
)

/*
	controller.go 用于设计 反射注入的处理器
	aurora的参数注入规则:
	URL  参数永远排列在参数列表的最前,
	GET  参数跟随在URL之后
	POST 排序在最后

	基于go反射的特点，处理器传递参数和前端传递参数的名称没有任何关系，只和顺序有关(调用服务器接口就和调用函数传参一样，需要给对对应类型)

*/
// control 用于存储在服务器启动之前注册的接口信息，需要在加载完配置项之后进行统一注册
type controlInfo struct {
	path       string
	control    Controller
	middleware []Middleware
}

type controller struct {
	*Aurora
	//解析的文件部分
	File *MultipartFile
	//上下文数据
	ctx Ctx
	p   *Proxy
	//路径参数,按顺序依次
	UrlVariable []string
	RESTFul     map[string]interface{}
	//入参参数个数
	InNum int
	//返回值个数
	OutNum int
	//按顺序存储每个入参的反射实例
	InvokeValues []reflect.Value
	//参数赋值序列表
	Args []string
	//可赋值参数索引序列
	AssignmentIndex []int
	//返回参数实例
	ReturnValues []reflect.Value
	//将被调用的函数,注册阶段已经被构建成为反射类型
	Fun     reflect.Value
	FunType reflect.Type
}

// InitArgs 初始化参数信息，注册函数阶段调用
func (c *controller) InitArgs() {
	c.InNum = c.FunType.NumIn()
	c.OutNum = c.FunType.NumOut()
	c.AssignmentIndex = make([]int, 0)
	//初始化参数列表
	if c.InNum > 0 {
		c.InvokeValues = make([]reflect.Value, c.InNum)
		c.Args = make([]string, c.InNum)
	}
	for i := 0; i < c.InNum; i++ {
		arguments := c.FunType.In(i)
		value := reflect.New(arguments).Elem()
		//初始化参数期间对参数列表进行标记，以便匹配参数顺序,此处主要是处理存在web请求体或者响应体的位置
		key := arguments.String()
		if _, b := c.Aurora.intrinsic[key]; b {
			c.Args[i] = key
			c.InvokeValues[i] = value
			continue
		}
		//对非内部参数进行 字段校验 存在为导出字段需要更改
		if arguments.Kind() == reflect.Struct || arguments.Kind() == reflect.Ptr {
			if !checkArguments(value) {
				//检查存在 未导出字段
				log.Fatalln("The index: ", i, " parameter is checked to exist as an export field, please check the field permission")
			}
		}
		c.InvokeValues[i] = value
		//初始化可赋值参数序列，存储可赋值的索引
		c.AssignmentIndex = append(c.AssignmentIndex, i)
	}
}

// checkArguments 校验接口入参 参数所有字段是否为导出字段
// 找要有一个是非导出字段则返回 false
func checkArguments(s reflect.Value) bool {
	var v reflect.Value
	if s.Kind() != reflect.Struct && s.Kind() != reflect.Ptr {
		return true
	}
	//如果入参是指针
	if s.Kind() == reflect.Ptr {
		//校验入参 此刻的指针数据是未初始化情况 需要分配一个值来进行校验,分配的值仅用于校验
		elem := reflect.New(s.Type().Elem()).Elem()
		return checkArguments(elem)
	} else {
		v = s
	}
	st := v.Type()
	for i := 0; i < st.NumField(); i++ {
		//兼容1.16 取消校验
		//field := st.Field(i)
		// 校验当前结构体的字段是否是导出状态
		//if !field.IsExported() {
		//	return false
		//}
		//对该字段进行递归检查
		if !checkArguments(v.Field(i)) {
			return false
		}
	}
	return true
}

// invoke 接口调用
func (c *controller) invoke() []reflect.Value {
	//before
	r := c.Fun.Call(c.InvokeValues)
	//after
	return r
}

//入参解析
func (c *controller) analysisInput(request *http.Request, response http.ResponseWriter, ctx Ctx) {
	// var values []string 用于接收 参数列表，该列表顺序规则为(rest full URL参数永远放在最前):
	// values:   [rest ful路径参数,GET 请求参数,POST请求体参数]
	var values []string
	//根据 请求类型初始化 values 列表
	switch request.Method {
	case http.MethodGet:
		values = getRequest(request, c)
	case http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodHead:
		values = getRequest(request, c)
		post := postRequest(request, c)
		values = append(values, post...)
	}
	if values == nil {
		//如果 values 在前面的请求类型中均未被初始化，则初始化 values为空零元素切片,以防下面出现空指针错误
		values = make([]string, 0)
	}
	//如果get请求参数个数少于处理函数入参数量，则按照get参数个数初始化，剩余部分判断为空，或是web请求提参数
	l := 0
	//校验 传递参数的数量和 可输入的数量
	if len(c.AssignmentIndex) <= len(values) {
		//传递参数的数量大于 处理函数可赋值数量，则默认丢弃多余部分的参数，以可赋值参数长度为主
		l = len(c.AssignmentIndex)
	} else {
		//此情况 传入的参数小于可赋值参数数量，把传入的参数按照可赋值索引依次赋值
		l = len(values)
	}
	//初始化参数列表，如果 values 为零个元素 则 不会给Args入参 进行初始化
	for i := 0; i < l; i++ {
		v := ""
		if strings.Contains(values[i], "=") {
			vs := strings.Split(values[i], "=")
			v = vs[1]
		} else {
			v = values[i]
		}
		assig := c.AssignmentIndex[i]
		c.Args[assig] = v
	}

	//开始初始化参数注入，Args中的参数没有被初始化 依然为 "" 空字符串，则在初始化的时候默认为 零值
	for i := 0; i < c.InNum; i++ {
		v := c.Args[i]
		if v == "" {
			continue
		}
		json := jsoniter.ConfigCompatibleWithStandardLibrary
		var data interface{}
		var err error
		if vr, b := c.p.Aurora.intrinsic[v]; b {
			prama := vr(c.p)
			pv := reflect.ValueOf(prama)
			if !pv.Type().AssignableTo(c.InvokeValues[i].Type()) {
				panic("The required type is'" + c.InvokeValues[i].Type().String() + "' The provided type is '" + pv.Type().String() + "'" +
					",Custom system parameter initialization error, please check whether the type returned by the constructor matches the type required by the processor")
			}
			c.InvokeValues[i] = reflect.ValueOf(prama)
			continue
		}
		if request.Method != http.MethodGet {
			err = json.Unmarshal([]byte(v), &data)
			if err != nil {
				panic("The json parameter decoding failed, please check whether the json data format is correct.error:" + err.Error())
			}
		} else {
			switch c.InvokeValues[i].Kind() {
			case reflect.Map, reflect.Struct, reflect.Interface, reflect.Ptr:
				if c.InvokeValues[i].Kind() == reflect.Ptr {
					kind := c.InvokeValues[i].Type().Elem().Kind()
					if !(kind == reflect.Map || kind == reflect.Struct) {
						data = v
						break
					}
				}
				query := request.URL.Query()
				if c.RESTFul == nil {
					c.RESTFul = map[string]interface{}{}
				}
				for k, v := range query {
					c.RESTFul[k] = v[0]
				}
				data = c.RESTFul
				//使用结构体或者map进行解析 在对应的参数位置应该多添加一个占位符号，以确保后面存在的参数能够正确被初始化复制，此处需要 在 i位置对 Args 添加一个占位
				s := c.Args[:i]
				e := c.Args[i:]
				c.Args = make([]string, 0)
				c.Args = append(c.Args, s...)
				c.Args = append(c.Args, "")
				c.Args = append(c.Args, e...)
			case reflect.Int, reflect.Float64, reflect.Bool, reflect.String, reflect.Float32, reflect.Int32:
				data = v
			}
		}
		err = utils.StarAssignment(c.InvokeValues[i], data)
		if err != nil {
			panic(err)
		}
	}

}

func getRequest(request *http.Request, c *controller) []string {
	values := make([]string, 0)
	url := request.RequestURI
	//解析存在get参数
	if index := strings.Index(url, "?"); index != -1 {
		url = url[index+1:]
		if c.UrlVariable != nil {
			//如果存在路径参数,我们把路径参数附加在 get参数之后
			values = c.UrlVariable
		}
		value := strings.Split(url, "&")
		values = append(values, value...)
	} else {
		if c.UrlVariable != nil {
			//如果存在路径参数,我们把路径参数附加在 get参数之后
			values = c.UrlVariable
		}
	}
	return values
}

func postRequest(request *http.Request, c *controller) []string {
	values := make([]string, 0)
	//处理文件上传处理 该处理操作在 中间件阶段可能被执行，两种情况同时出现的情况未测试，可能出现bug
	request.ParseMultipartForm(c.p.MaxMultipartMemory)
	form := request.MultipartForm
	if form != nil {
		if form.File != nil {
			//封装解析好的 文件部分
			c.File = &MultipartFile{File: form.File}
		}
		if form.Value != nil {
			// 2022-5-20 更新 多文本混合上传方式
			for _, v := range form.Value {
				vlen := len(v)
				if vlen == 0 {
					continue
				}
				values = append(values, v[0])
			}
			return values
		}
	}
	//非文件上传处理,可能存在bug
	if request.Body != nil {
		all, err := ioutil.ReadAll(request.Body)
		if err != nil {
			//待处理
		}
		//确保读取到内容
		if all != nil && len(all) > 0 {
			values = append(values, string(all))
		}
	}
	return values
}

// Control 初始化装配结构体依赖 control 参数必须是指针
func (a *Aurora) control(control Controller) {
	value, err := checkControl(control)
	if err != nil {
		a.Panic(err)
	}
	if a.controllers == nil {
		a.controllers = make([]*reflect.Value, 0)
	}
	a.controllers = append(a.controllers, value)
	// 把处理器注册进 ioc , 默认为类型名称
	tf := reflect.TypeOf(control)
	err = a.component.putIn(tf.String(), control)
	if err != nil {
		a.Panic(err)
	}
}

// checkControl 校验处理器的规范形式
func checkControl(control Controller) (*reflect.Value, error) {
	v := reflect.ValueOf(control)
	//指针类型校验
	if v.Kind() != reflect.Ptr {
		return nil, errors.New("'" + v.Type().String() + "' not pointer, requires a pointer parameter")
	}
	//空指针校验
	if v.IsNil() {
		return nil, errors.New("null pointer")
	}
	//指针类型结构体校验
	if v.Elem().Kind() != reflect.Struct {
		return nil, errors.New("requires a struct type")
	}
	return &v, nil
}

const ref = "ref"

// dependencyInjection Control 依赖加载
func (a *Aurora) dependencyInjection() {
	if a.controllers == nil {
		return
	}
	l := len(a.controllers)
	for i := 0; i < l; i++ {
		control := *a.controllers[i]
		if control.Kind() == reflect.Ptr {
			control = control.Elem()
		}
		for j := 0; j < control.NumField(); j++ {
			field := control.Type().Field(j)
			//t := field.Type.Kind()
			//if t == reflect.Ptr {
			//	t = field.Type.Elem().Kind()
			//}
			//if t != reflect.Struct {
			//	continue
			//}
			//查找结构体字段的tag 找到 ref属性并初始化 Get 兼容1.16
			//if r, b := field.Tag.Lookup(ref); b {
			//	//if r == "" {
			//	//	//暂时选择警告的方式
			//	//	a.Warn("ref tag value is ''")
			//	//	continue
			//	//}
			//	//校验字段是否可操作 兼容1.16 注释掉 检查字段操作
			//	//if !field.IsExported() {
			//	//	a.Panic(control.Type().String(), " '", field.Name, "' field .Not an exported field, ref attribute injection failed")
			//	//}
			//	value := a.component.get(r)
			//	if value == nil {
			//		//若是 容器中未找到 ref 则初始化失败 结束程序
			//		a.Panic("not found ref: " + r)
			//	}
			//	//容器中找到了 ref 属性的 参数，把该参数赋值给 给定的处理器字段
			//	if err := utils.Injection(control.Field(j), *value); err != nil {
			//		//初始化赋值阶段错误
			//		a.Panic(err.Error())
			//	}
			//} else {
			//	//尝试获取 同类型的依赖
			//	value := a.component.get(control.Field(j).Type().String())
			//	if value != nil {
			//		if err := utils.Injection(control.Field(j), *value); err != nil {
			//			//初始化赋值阶段错误
			//			a.Panic(err.Error())
			//		}
			//	}
			//}
			//查询 value 属性 读取config中的基本属性
			if v, b := field.Tag.Lookup("value"); b {
				if v == "" {
					a.Warn("value tag value is ''")
					continue
				}
				get := a.config.Get(v)
				if get == nil {
					//如果查找结果大小等于0 则表示不存在
					continue
				}
				//把查询到的 value 初始化给指定字段
				err := utils.StarAssignment(control.Field(j), get)
				if err != nil {
					a.Panic(err.Error())
				}
			}
		}
	}
}
