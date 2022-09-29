package aurora

import (
	"gitee.com/aurora-engine/aurora/route"
	"gitee.com/aurora-engine/aurora/utils"
	"reflect"
)

type Constructor func(*route.Proxy) interface{}

func (engine *Engine) SysVariable(v interface{}, constructor route.Constructor) {
	if v == nil || constructor == nil {
		return
	}
	rt := reflect.ValueOf(v)
	if engine.intrinsic == nil {
		engine.intrinsic = make(map[string]route.Constructor)
	}
	key := utils.BaseTypeKey(rt)
	engine.intrinsic[key] = constructor
}

// 系统变量

func req(proxy *route.Proxy) interface{} {
	return proxy.Req
}

func rew(proxy *route.Proxy) interface{} {
	return proxy.Rew
}

func ctx(proxy *route.Proxy) interface{} {
	return proxy.Context
}

func file(proxy *route.Proxy) interface{} {
	return proxy.File
}
