package route

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

func ErrorMsg(err error, msg ...string) {
	if err != nil {
		if msg == nil {
			msg = []string{"Error"}
		}
		emsg := fmt.Errorf("%s : %s", strings.Join(msg, ""), err.Error())
		panic(emsg)
	}
}

// Catch 错误捕捉器的存储上要进行封装
type Catch struct {
	in  []reflect.Value
	fun reflect.Value
}

func (c *Catch) invoke(err reflect.Value) []reflect.Value {
	c.in[0] = err
	return c.fun.Call(c.in)
}

func (router *Router) registerErrorCatch(err any) {
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
	c := Catch{
		in: []reflect.Value{
			reflect.New(in).Elem(),
		},
		fun: reflect.ValueOf(err),
	}
	if router.Catches == nil {
		router.Catches = map[reflect.Type]Catch{in: c}
		return
	}
	router.Catches[in] = c
}
