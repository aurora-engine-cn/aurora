package route

import (
	"fmt"
	"gitee.com/aurora-engine/aurora/core"
	"gitee.com/aurora-engine/aurora/web"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"reflect"
	"strconv"
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

type Controller struct {
	*Proxy
	Constraints     map[string]web.Verify
	Context         web.Context     //上下文数据
	UrlVariable     []string        //路径参数,按顺序依次
	RESTFul         map[string]any  // K/V 路径参数 控制器解析中，请求参数参数名尽量不要和 RESTFul 参数同名，否则会覆盖
	InNum           int             //处理器入参参数个数
	OutNum          int             //处理器返回值个数
	InvokeValues    []reflect.Value // InvokeValues存储的是控制器传递参数的序列 按顺序存储每个入参的反射实例
	Args            []string        //参数赋值序列表，主要存储请求参数的只值
	AssignmentIndex []int           //AssignmentIndex 可赋值参数索引序列，可赋值参数序列是存储了系统内部参数之外的请求参数所在 InvokeValues 参数序列中的索引位置。
	ReturnValues    []reflect.Value //返回参数实例
	Fun             reflect.Value   //将被调用的函数,注册阶段已经被构建成为反射类型
	FunType         reflect.Type
	Intrinsic       map[string]web.Variate // 自定赋值参数列表(系统参数配置)
}

// InitArgs 初始化参数信息，注册函数阶段调用
// 完成对 InvokeValues 控制器参数的初始化(未赋值状态)
// 完成对应的 AssignmentIndex 可赋值参数序列初始化
func (control *Controller) InitArgs() {
	control.InNum = control.FunType.NumIn()
	control.OutNum = control.FunType.NumOut()
	control.AssignmentIndex = make([]int, 0)
	//初始化参数列表
	if control.InNum > 0 {
		control.InvokeValues = make([]reflect.Value, control.InNum)
		control.Args = make([]string, control.InNum)
	}
	for i := 0; i < control.InNum; i++ {
		arguments := control.FunType.In(i)
		value := reflect.New(arguments).Elem()
		//初始化参数期间对参数列表进行标记，以便匹配参数顺序,此处主要是处理存在web请求体或者响应体的位置
		key := core.BaseTypeKey(value)
		if _, b := control.Intrinsic[key]; b {
			control.Args[i] = key
			control.InvokeValues[i] = value
			continue
		}
		control.InvokeValues[i] = value
		//初始化可赋值参数序列，存储可赋值的索引
		control.AssignmentIndex = append(control.AssignmentIndex, i)
	}
}

// invoke 接口调用
func (control *Controller) invoke() []reflect.Value {
	//before
	// 结构体参数约束校验
	err := control.checkConstrain()
	if err != nil {
		panic(err)
	}
	r := control.Fun.Call(control.InvokeValues)
	//after
	return r
}

// 入参解析
func (control *Controller) analysisInput(request *http.Request) {
	// var values []string 用于接收 参数列表，该列表顺序规则为(rest full URL参数永远放在最前):
	// values:   [rest ful路径参数,GET 请求参数,POST请求体参数]
	var values []string
	//根据 请求类型初始化 values 列表
	switch request.Method {
	case http.MethodGet:
		values = GetRequest(request, control)
	case http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodHead:
		values = GetRequest(request, control)
		postValue := PostRequest(request, control)
		values = append(values, postValue...)
	}
	if values == nil {
		//如果 values 在前面的请求类型中均未被初始化，则初始化 values为空零元素切片,以防下面出现空指针错误
		values = make([]string, 0)
	}
	//如果get请求参数个数少于处理函数入参数量，则按照get参数个数初始化，剩余部分判断为空，或是web请求提参数
	l := 0
	//校验 传递参数的数量和 可输入的数量
	if len(control.AssignmentIndex) <= len(values) {
		//传递参数的数量大于 处理函数可赋值数量，则默认丢弃多余部分的参数，以可赋值参数长度为主
		l = len(control.AssignmentIndex)
	} else {
		//此情况 传入的参数小于可赋值参数数量,传入的参数按照可赋值索引依次赋值
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
		assign := control.AssignmentIndex[i]
		control.Args[assign] = v
	}

	//开始初始化参数注入，Args中的参数没有被初始化 依然为 "" 空字符串，则在初始化的时候默认为 零值
	for i := 0; i < control.InNum; i++ {
		v := control.Args[i]
		if v == "" {
			continue
		}
		var data interface{}
		var err error
		if vr, b := control.Intrinsic[v]; b {
			parameter := vr(control.Context)
			pv := reflect.ValueOf(parameter)
			if !pv.Type().AssignableTo(control.InvokeValues[i].Type()) {
				panic("The required type is'" + control.InvokeValues[i].Type().String() + "' The provided type is '" + pv.Type().String() + "'" +
					",Custom system parameter initialization error, please check whether the type returned by the constructor matches the type required by the processor")
			}
			control.InvokeValues[i] = reflect.ValueOf(parameter)
			continue
		}
		if request.Method != http.MethodGet {
			err = jsoniter.Unmarshal([]byte(v), &data)
			ErrorMsg(err, "The json parameter decoding failed, please check whether the json data format is correct.error:")
		} else {
			switch control.InvokeValues[i].Kind() {
			case reflect.Map, reflect.Struct, reflect.Interface, reflect.Ptr:
				if control.InvokeValues[i].Kind() == reflect.Ptr {
					kind := control.InvokeValues[i].Type().Elem().Kind()
					if !(kind == reflect.Map || kind == reflect.Struct) {
						data = v
						break
					}
				}
				query := request.URL.Query()
				if control.RESTFul == nil {
					control.RESTFul = map[string]any{}
				}
				for k, v := range query {
					control.RESTFul[k] = v[0]
				}
				data = control.RESTFul
				//使用结构体或者map进行解析 在对应的参数位置应该多添加一个占位符号，以确保后面存在的参数能够正确被初始化复制，此处需要 在 i位置对 Args 添加一个占位
				s := control.Args[:i]
				e := control.Args[i:]
				control.Args = make([]string, 0)
				control.Args = append(control.Args, s...)
				control.Args = append(control.Args, "")
				control.Args = append(control.Args, e...)
			case reflect.Int, reflect.Float64, reflect.Bool, reflect.String, reflect.Float32, reflect.Int32:
				data = v
			}
		}
		err = core.StarAssignment(control.InvokeValues[i], data)
		ErrorMsg(err)
	}

}

func GetRequest(request *http.Request, control *Controller) []string {
	values := make([]string, 0)
	url := request.RequestURI
	//解析存在get参数
	if index := strings.Index(url, "?"); index != -1 {
		url = url[index+1:]
		if control.UrlVariable != nil {
			//如果存在路径参数,我们把路径参数附加在 get参数之后
			values = control.UrlVariable
		}
		value := strings.Split(url, "&")
		values = append(values, value...)
	} else {
		if control.UrlVariable != nil {
			//如果存在路径参数,我们把路径参数附加在 get参数之后
			values = control.UrlVariable
		}
	}
	return values
}

func PostRequest(request *http.Request, control *Controller) []string {
	values := make([]string, 0)
	//处理文件上传处理 该处理操作在 中间件阶段可能被执行，两种情况同时出现的情况未测试，可能出现bug
	request.ParseMultipartForm(control.MaxMultipartMemory)
	form := request.MultipartForm
	if form != nil {
		if form.File != nil {
			//封装解析好的 文件部分
			control.Context[web.AuroraMultipartFile] = &web.MultipartFile{File: form.File}
		}
		if form.Value != nil {
			// 2022-5-20 更新 多文本混合上传方式
			for _, v := range form.Value {
				length := len(v)
				if length == 0 {
					continue
				}
				values = append(values, v[0])
			}
			return values
		}
	}
	//非文件上传处理,可能存在bug
	if request.Body != nil {
		all, err := io.ReadAll(request.Body)
		if err != nil {
			//待处理
			panic(err)
		}
		//确保读取到内容
		if all != nil && len(all) > 0 {
			values = append(values, string(all))
		}
	}
	return values
}

// 检查结构体参数中的约束是否满足对应检查
func (control *Controller) checkConstrain() error {
	for i := 0; i < len(control.InvokeValues); i++ {
		if ok, err := control.check(control.InvokeValues[i]); !ok {
			// 待优化消息提示具体到函数名称
			return err
		}
	}
	return nil
}

func (control *Controller) check(value reflect.Value) (bool, error) {
	if value.Kind() == reflect.Ptr {
		return control.check(value.Elem())
	}
	if value.Kind() == reflect.Struct {
		// 校验各个 字段的 tar
		fields := value.NumField()
		for i := 0; i < fields; i++ {
			field := value.Type().Field(i)
			tag := field.Tag

			// 检查 empty 空值校验,empty 为false 表示该字段不能为空 true 表示字段可以为空
			empty, b := tag.Lookup("empty")
			if b && empty != "" {
				parseBool, err := strconv.ParseBool(empty)
				ErrorMsg(err, "tag:empty '"+empty+"' value could not be parsed")
				if value.Field(i).IsZero() && !parseBool {
					// 校验不通过
					return false, fmt.Errorf("the '%s' controller in '%s.%s' constraint check failed", control.FunType.String(), control.InvokeValues[i].Type().String(), field.Name)
				}
			}

			// 检查默认值 constraint 的值应该是一个字符串,并且使用分号隔开
			// constraint 用于调用自定义约束
			constraint, b := tag.Lookup("constraint")
			if b && constraint != "" {
				tags := strings.Split(constraint, ";")
				for j := 0; j < len(tags); j++ {
					tagKey := tags[j]
					tagKey = strings.TrimSpace(tagKey)
					if ConstraintFunc, flag := control.Constraints[tagKey]; flag {
						err := ConstraintFunc(value.Field(i).Interface())
						if err != nil {
							return false, fmt.Errorf("the '%s' controller in '%s.%s' constraint check failed. %s", control.FunType.String(), control.InvokeValues[i].Type().String(), field.Name, err.Error())
						}
					}
				}
			}

		}
	}
	return true, nil
}
