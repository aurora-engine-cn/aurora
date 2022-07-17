package aurora

import "reflect"

type Constructor func(*Proxy) interface{}

func (a *Aurora) SysVariable(v interface{}, constructor Constructor) {
	if v == nil || constructor == nil {
		return
	}
	rt := reflect.TypeOf(v)
	if a.intrinsic == nil {
		a.intrinsic = make(map[string]Constructor)
	}
	a.intrinsic[rt.String()] = constructor
}

// 系统变量

func req(proxy *Proxy) interface{} {
	return proxy.Req
}

func rew(proxy *Proxy) interface{} {
	return proxy.Rew
}

func ctx(proxy *Proxy) interface{} {
	return proxy.Ctx
}

func file(proxy *Proxy) interface{} {
	return proxy.File
}
