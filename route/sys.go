package route

type Constructor func(*Proxy) interface{}

//func (engine *Engine) SysVariable(v interface{}, constructor Constructor) {
//	if v == nil || constructor == nil {
//		return
//	}
//	rt := reflect.ValueOf(v)
//	if engine.intrinsic == nil {
//		engine.intrinsic = make(map[string]Constructor)
//	}
//	key := utils.BaseTypeKey(rt)
//	engine.intrinsic[key] = constructor
//}

// 系统变量

func req(proxy *Proxy) interface{} {
	return proxy.Req
}

func rew(proxy *Proxy) interface{} {
	return proxy.Rew
}

func ctx(proxy *Proxy) interface{} {
	return proxy.Context
}

func file(proxy *Proxy) interface{} {
	return proxy.File
}
