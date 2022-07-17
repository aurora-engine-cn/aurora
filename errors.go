package aurora

import (
	"encoding/json"
	"log"
	"reflect"
	"time"
)

/*
	错误处理
*/

type ErrorResponse struct {
	UrlPath      string `json:"url"`
	Status       int    `json:"code"`
	ErrorMessage string `json:"error"`
	Time         string `json:"time"`
}

func newErrorResponse(path, message string, status int) string {
	now := time.Now().Format("2006/01/02 15:04:05")
	msg := ErrorResponse{
		UrlPath:      path,
		Status:       status,
		ErrorMessage: message,
		Time:         now,
	}
	marshal, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return string(marshal)
}

// ArgsAnalysisError 参数解析错误
type ArgsAnalysisError struct {
	s string
}

func (a ArgsAnalysisError) Error() string {
	return a.s
}

// Error 错误类型  类型设计 是一个函数 接收一个 实现了 error 接口的参数
type Error = interface{}

//错误捕捉器的存储上要进行封装
type catch struct {
	in  []reflect.Value
	fun reflect.Value
}

func (c *catch) invoke(err reflect.Value) []reflect.Value {
	c.in[0] = err
	return c.fun.Call(c.in)
}

func (a *Aurora) Catch(err Error) {
	a.router.Catch(err)
}

func (r *route) registerErrorCatch(err Error) {
	if err == nil {
		return
	}
	of := reflect.TypeOf(err)
	if of.Kind() != reflect.Func && of.NumIn() != 1 {
		return
	}
	//校验入参是否为 error
	in := of.In(0)
	e := new(error)
	et := reflect.TypeOf(e).Elem()
	if !in.Implements(et) {
		log.Panic(of.Name() + " not is errors!")
		return
	}
	c := catch{
		in: []reflect.Value{
			reflect.New(in).Elem(),
		},
		fun: reflect.ValueOf(err),
	}
	if r.catch == nil {
		r.catch = map[reflect.Type]catch{in: c}
		return
	}
	r.catch[in] = c
}
