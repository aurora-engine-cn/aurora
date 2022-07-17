package aurora

import (
	"errors"
	"github.com/aurora-go/aurora/utils"
	"reflect"
	"sync"
)

/*
	管理集成解耦关系 进行中
	Ioc:
	1.整个容器中 id是唯一的

*/

func newIoc() *ioc {
	return &ioc{
		id:    &sync.Map{},
		first: make(map[string]*reflect.Value),
	}
}

type ioc struct {
	kv    map[string]interface{}    // 初始缓存容器
	first map[string]*reflect.Value //一级,存放有部分属性未初始化的容器
	id    *sync.Map                 //主缓存容器
}

// put 向容器中添加一个 注册，相同的id直接覆盖
func (i *ioc) put(ref string, value interface{}) {
	if value == nil {
		return
	}
	v := reflect.ValueOf(value)
	i.id.Store(ref, &v)
}

// get 从容器中获取引用,存指针的目的在于方便校验 nil
func (i *ioc) get(ref string) *reflect.Value {
	load, ok := i.id.Load(ref)
	if !ok {
		return nil
	}
	v := load.(*reflect.Value)
	return v
}

// putIn 想容器装载属性,不允许覆盖已经存在的 属性实例，主要提供对外使用.
// 它允许覆盖 初级缓存中的配置项
func (i *ioc) putIn(ref string, value interface{}) error {
	if _, b := i.id.Load(ref); b {
		//不能注册 已存在的 key
		return errors.New("'" + ref + "' id already exists, repeated registration failed")
	}
	if i.kv == nil {
		i.kv = map[string]interface{}{
			ref: value,
		}
		return nil
	}

	if _, b := i.kv[ref]; b {
		return errors.New("'" + ref + "' id already exists, repeated registration failed")
	}
	i.kv[ref] = value
	return nil
}

// update 用于更新容器中的实例,如果存在就更新，否则不做任何动作，返回一个bool 用于判断跟新结果，如果更新了 则需要重新对 Aurora 的 控制器进行一次赋值更新
func (i *ioc) update(ref string, value interface{}) bool {
	if _, b := i.id.Load(ref); !b {
		return false
	}
	v := reflect.ValueOf(value)
	i.id.Store(ref, &v)
	return true
}

// start 启动容器
// 启动容器将 kv装在的属性进行初始化
func (i *ioc) start() error {
	if i.kv == nil {
		return nil
	}
	return i.run()
}

func (i *ioc) run() error {
	// 初始 缓存 加载
	for k, v := range i.kv {
		value := reflect.ValueOf(v)
		if err := i.dependence(k, value); err != nil {
			return err
		}
	}
	// 一级缓存 加载
	for k, v := range i.first {
		if err := i.dependence(k, *v); err != nil {
			return err
		}
	}
	return nil
}

// dependence
func (i *ioc) dependence(ref string, value reflect.Value) error {
	// values 主要用来操作结构体字段
	var values reflect.Value
	// 检擦 value 是否为指针
	if value.Kind() == reflect.Ptr {
		// 我们需要 操作指向的值进行初始化
		values = value.Elem()
		//需要检擦 指向的值是否为结构体，存在双重指针或者多级指针的 视为无效属性
		if values.Kind() == reflect.Ptr {
			return errors.New("invalid parameter, there is a double pointer or a multi-level pointer")
		}
	} else {
		values = value
	}
	//开始扫描该value是否有需要装配的属性
	for j := 0; j < values.NumField(); j++ {
		field := values.Type().Field(j)
		t := field.Type.Kind()
		if t == reflect.Ptr {
			t = field.Type.Elem().Kind()
		}
		if t != reflect.Struct {
			continue
		}
		// 校验容器中的组件属性是否被初始化过，未初始化则交由容器初始化
		if !values.Field(j).IsZero() {
			continue
		}

		// r 是我们需要去 id主容器 中查找的依赖项
		if r, b := field.Tag.Lookup("ref"); b {
			//if r == "" {
			//	//检测 ref 属性是否为空字符串,为空则跳过
			//	continue
			//}
			////检查是否导出,被操纵字段必须是可导出的
			//if !field.IsExported() {
			//	//不可导出出字段无法赋值
			//	return errors.New("ref attribute on non-exported field, cannot complete initialization")
			//}
			var va *reflect.Value
			//开始查询 tag 的 引用id 是否在主容器中
			load, ok := i.id.Load(r)
			if !ok {
				// 主容器中没有找到 tag 引用，尝试在 一缓存容器中查找
				v, o := i.first[r]
				if !o {
					// 一级缓存容器 也无法查询到
					//此刻需要 把当前正在初始化的 实例放到 第一级缓存容器中 表示 当前的实例已经存在，然后去 初始化未找到的引用
					// 存放到一级缓存 之前需要判定当前 初始化的实例 是否已经在第一缓存中，这里判断的目的主要是校验 完成初级缓存加载后，继续进行第二次 缓存加载依然无法找到指定的 tag 引用，这个情况下找不到的 引用只可能是不存在于所有容器中，在第二次缓存加载走到这里就是错误的
					_, f := i.first[ref]
					if f {
						// 第一次缓存加载 ，此处必定不会执行，若是第二次缓存加载 ，并且没有找到指定的 ref 必定走到此处 将返回错误
						return errors.New(r + " Reference instance not found")
					}
					i.first[ref] = &value
					// 存储完成后 删除原来 kv中的该实例 ,以防下次重复
					delete(i.kv, ref)
					return i.run()
				}
				va = v
			}
			if load != nil {
				va = load.(*reflect.Value)
			}
			//如果找到了ref 校验是否可以赋值给该字段
			if va == nil {
				return errors.New("initialization failed")
			}
			if !va.Type().AssignableTo(field.Type) {
				return errors.New("cannot assign type mismatch")
			}
			//进行 赋值
			if err := utils.Injection(values.Field(j), *va); err != nil {
				return err
			}
		} else {
			//尝试 通过类型匹配寻找 ( 结构体同类型匹配查找 )，此处可能存在bug 如果容器的该字段属性是存在初始化好的，此处会覆盖掉原来的赋值
			var va *reflect.Value
			r = field.Type.String()
			//开始查询 tag 的 引用id 是否在主容器中
			load, ok := i.id.Load(r)
			if !ok {
				// 主容器中没有找到 tag 引用，尝试在 一缓存容器中查找
				v, o := i.first[r]
				if !o {
					// 一级缓存容器 也无法查询到
					//此刻需要 把当前正在初始化的 实例放到 第一级缓存容器中 表示 当前的实例已经存在，然后去 初始化未找到的引用
					// 存放到一级缓存 之前需要判定当前 初始化的实例 是否已经在第一缓存中，这里判断的目的主要是校验 完成初级缓存加载后，继续进行第二次 缓存加载依然无法找到指定的 tag 引用，这个情况下找不到的 引用只可能是不存在于所有容器中，在第二次缓存加载走到这里就是错误的
					_, f := i.first[ref]
					if f {
						// 第一次缓存加载 ，此处必定不会执行，若是第二次缓存加载 ，并且没有找到指定的 ref 必定走到此处 将返回错误 此处作为适配 结构体类型加载，找不到视为没有不影响初始化 跳过这个字段
						//return errors.New(r + " Reference instance not found")
						// 当前的 结构属性未初始化完毕
						continue
					}
					i.first[ref] = &value
					// 存储完成后 删除原来 kv中的该实例 ,以防下次重复
					delete(i.kv, ref)
					return i.run()
				}
				va = v
			}
			if load != nil {
				va = load.(*reflect.Value)
			}
			//如果找到了ref 校验是否可以赋值给该字段
			if va == nil {
				return errors.New("initialization failed")
			}
			if !va.Type().AssignableTo(field.Type) {
				return errors.New("cannot assign type mismatch")
			}
			//进行 赋值
			if err := utils.Injection(values.Field(j), *va); err != nil {
				return err
			}
		}
	}
	if _, b := i.first[ref]; b {
		//完成初始化后，如果该实例存在于一级缓存中 我们需要把它从一级缓存中删除
		delete(i.first, ref)
	}
	// 重新存到主缓存中
	i.id.Store(ref, &value)
	delete(i.kv, ref)
	return nil
}
