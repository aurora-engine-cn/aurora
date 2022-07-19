package aurora

import (
	jsoniter "github.com/json-iterator/go"
	"log"
	"net/http"
	"reflect"
)

// Aurora 全局错误 panic 处理
func errRecover(proxy *Proxy) {
	rew := proxy.Rew
	if v := recover(); v != nil {
		var msg string
		switch v.(type) {
		case error:
			msg = v.(error).Error()
		case string:
			msg = v.(string)
		default:
			marshal, err := jsoniter.Marshal(v)
			if err != nil {
				msg = err.Error()
			}
			msg = string(marshal)
		}
		proxy.Error(msg)
		http.Error(rew, msg, 500)
		return
	}
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
