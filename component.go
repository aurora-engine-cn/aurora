package aurora

import (
	"gitee.com/aurora-engine/aurora/core"
	"reflect"
)

// Component 命名组件
type Component map[string]interface{}

// Constructors 用于加载 匿名组件的构造器类型
// Aurora 会执行这个函数 并把得到的变量加载到 ioc 容器中
type Constructors func() interface{}

// StartIoc 启动容器
func (engine *Engine) ioc() {
	engine.Info("start component-dependent assembly")

	//加载uses配置项，配置项中可能存在加载ioc配置
	if engine.options != nil {
		for _, useOption := range engine.options {
			useOption(engine)
		}
	}
	// 加载 构造器 build 到 ioc 容器
	if engine.build != nil {
		for _, constructor := range engine.build {
			// 执行构造 生成组件放入到 ioc中
			c := constructor()
			err := engine.space.Put("", c)
			if err != nil {
				panic(err)
			}
		}
	}
	if engine.components != nil {
		for _, component := range engine.components {
			for k, v := range component {
				if err := engine.space.Put(k, v); err != nil {
					ErrorMsg(err)
				}
			}
		}
	}
	// 清空
	engine.components = nil
	//启动容器 ,给容器中的组件进行依赖初始化,容器加载出错 结束运行
	err := engine.space.Start()
	if err != nil {
		ErrorMsg(err,"Container initialization failed")
	}
	engine.injection()
}

// injection  控制器依赖加载依赖加载,控制器的依赖加载实际在容器初始化阶段就已经完成
func (engine *Engine) injection() {

	// 获取容器中的主缓存
	Controllers:=engine.space.Cache()
	for _, c := range Controllers {
		control:=*c
		if control.Kind() == reflect.Ptr {
			control = control.Elem()
		}
		for j := 0; j < control.NumField(); j++ {
			field := control.Type().Field(j)
			//查询 value 属性 读取config中的基本属性
			if v, b := field.Tag.Lookup("value"); b && v!=""{
				get := engine.config.Get(v)
				if get == nil {
					//如果查找结果大小等于0 则表示不存在
					continue
				}
				//把查询到的 value 初始化给指定字段
				err := core.StarAssignment(control.Field(j), get)
				ErrorMsg(err)
			}
		}
	}
}