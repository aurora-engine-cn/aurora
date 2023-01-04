package container

// container.go 用于重构 ioc.go
// 从 go1.19 版本开始 container.go 只接受指针变量放入容器
// container 存在的问题，当前版本的容器仅支持启动程序运行一次，
// 为容器中的各个数据初始化赋值。虽然 Put Star等方法可以调用多次，但是container无法保证在协程中会发生什么样的问题
// 后续将考虑 将容器添加启动标识，让它只能进行一次启动。

import (
	"errors"
	"fmt"
	"reflect"
)

type name interface {
	Put(key string, value any) error
	Start() error
}

func NewSpace() *Space {
	return &Space{
		initializeCache: make(map[string]any),
		mainCache:       make(map[string]*reflect.Value),
		firstCache:      make(map[string]*reflect.Value),
	}
}

type Space struct {
	initializeCache    map[string]any            // 初始缓存容器
	firstCache         map[string]*reflect.Value // 一级缓存,存放有部分属性未初始化的容器
	mainCache          map[string]*reflect.Value // 主缓存容器
	defaultConstructor map[string]*reflect.Value // 默认构造函数

}

func (space *Space) Constructor(New ...any) {

}

// Put 向容器中存放依赖项
// key ：作为容器中的唯一标识，如果key为空字符串则会生成 value的默认key，默认key是包名全路径和字符串类型来表示
func (space *Space) Put(key string, value any) error {
	key, err := space.CacheCheck(key, value)
	if err != nil {
		return err
	}
	space.initializeCache[key] = value
	return nil
}

// Get 从容器中获取引用,存指针的目的在于方便校验 nil
func (space *Space) Get(ref string) *reflect.Value {
	load, ok := space.mainCache[ref]
	if !ok {
		return nil
	}
	return load
}

// Cache 获取主缓存
// 需要等待容器启动完成获取
func (space *Space) Cache() map[string]*reflect.Value {
	return space.mainCache
}

// Start 启动容器
// 启动容器将 kv装在的属性进行初始化
func (space *Space) Start() error {
	if space.initializeCache == nil || space.firstCache == nil || space.mainCache == nil {
		return errors.New("container not initialized")
	}
	return space.run()
}

// 运行容器
func (space *Space) run() error {
	// 初始 缓存 加载
	for k, v := range space.initializeCache {
		value := reflect.ValueOf(v)
		if err := space.dependence(k, value); err != nil {
			return err
		}
	}
	// 一级缓存 加载
	for k, v := range space.firstCache {
		if err := space.dependence(k, *v); err != nil {
			return err
		}
	}
	return nil
}

// dependence
// depKey : 当前属性的唯一依赖id
func (space *Space) dependence(depKey string, value reflect.Value) error {
	// values 主要用来操作结构体字段
	var values, fieldValue reflect.Value
	// 检擦 value 是否为指针
	switch value.Kind() {
	case reflect.Pointer:
		// 我们需要 操作指向的值进行初始化
		values = value.Elem()
		//需要检擦 指向的值是否为结构体，存在双重指针或者多级指针的 视为无效属性
		if values.Kind() == reflect.Pointer {
			return errors.New("invalid parameter, there is a double pointer or a multi-level pointer")
		}
	case reflect.Interface:

	case reflect.Struct:
		return nil
	default:
		return nil
	}
	//开始扫描该value是否有需要装配的属性
	for j := 0; j < values.NumField(); j++ {
		fieldValue = values.Field(j)
		// 获取到结构体字段的反射类型
		fieldType := values.Type().Field(j)
		// 获取字段类型的类别
		switch fieldType.Type.Kind() {
		case reflect.Pointer:
			kind := fieldType.Type.Elem().Kind()
			if !fieldValue.CanSet() || kind != reflect.Struct || !fieldValue.IsZero() {
				// 必须是可设置的
				// 校验容器中的组件属性是否被初始化过，未初始化则交由容器初始化
				// 检查字段需要满足 类别是结构体 并且是没有被初始化的
				continue
			}
		default:
			if !fieldValue.CanSet() || !fieldValue.IsZero() {
				continue
			}
		}
		var depValue *reflect.Value
		Key, check := DepKey(fieldType)
		//开始查询 tag 的 引用 id 是否在主容器中
		depValue, ok := space.mainCache[Key]
		if !ok {
			// 主容器中没有找到 tag 引用，尝试在 一缓存容器中查找
			v, find := space.firstCache[Key]
			if !find {
				// 一级缓存容器也无法查询到
				// 此刻需要把当前正在初始化的实例放到第一级缓存容器中表示当前的实例已经存在，然后去初始化未找到的引用(此处去初始化未找到的索引并不是直接去，而是通过下面的操作重新运行一次启动容器)
				// 存放到一级缓存 之前需要判定当前初始化的实例是否已经在第一缓存中,这里判断的目的主要是校验完成初级缓存加载后
				// 继续进行第二次缓存加载依然无法找到指定的 tag 引用，这个情况下找不到的引用只可能是不存在于容器中，在第二次缓存加载走到这里就是错误的
				_, is := space.firstCache[depKey]
				if is {
					msg := ""
					switch fieldType.Type.Kind() {
					case reflect.Pointer:
						msg = fmt.Sprintf("'%s-%s' '%s-%s' Reference instance not found \n", values.Type().PkgPath(), value.Type().String(), fieldValue.Type().Elem().PkgPath(), fieldValue.Type().String())
					case reflect.Interface:
						msg = fmt.Sprintf("'%s-%s' '%s-%s' Reference instance not found \n", values.Type().PkgPath(), value.Type().String(), fieldValue.Type().PkgPath(), fieldValue.Type().String())
					}
					// check 主要用于校验 这个字段是不是需要强制检验，强制检验主要是针对字段上面有 tag  ref属性的引用，ref引用找不到就会返回错误
					if check {
						// 第一次缓存加载 ，此处必定不会执行，若是第二次缓存加载 ，并且没有找到指定的 ref 必定走到此处 将返回错误
						return errors.New(msg)
					}
					fmt.Print(msg)
					// 跳过该属性的初始化
					continue
				}
				space.firstCache[depKey] = &value
				// 存储完成后 删除原来 kv中的该实例 ,以防下次重复
				delete(space.initializeCache, depKey)
				return space.run()
			}
			depValue = v
		}
		//如果找到了 ref 校验是否可以赋值给该字段
		if depValue == nil {
			return errors.New("initialization failed")
		}
		//进行 赋值
		if err := Injection(fieldValue, *depValue); err != nil {
			return err
		}
	}
	if _, b := space.firstCache[depKey]; b {
		//完成初始化后，如果该实例存在于一级缓存中 我们需要把它从一级缓存中删除
		delete(space.firstCache, depKey)
	}
	// 重新存到主缓存中
	space.mainCache[depKey] = &value
	delete(space.initializeCache, depKey)
	return nil
}

// CacheCheck 向容器添加缓存检查，返回检查错误或者生成key
// 容器只接受纸箱结构体的指针变量
func (space *Space) CacheCheck(key string, value any) (string, error) {
	if key == "" {
		key = TypeKey(value)
	}
	valueOf := reflect.ValueOf(value)
	if valueOf.Kind() != reflect.Ptr {
		return "", fmt.Errorf("'%s' is not a pointer type, please add a pointer type", key)
	}
	if _, b := space.mainCache[key]; b {
		//不能注册 已存在的 key
		return "", errors.New("'" + key + "' id already exists, repeated registration failed")
	}
	if space.initializeCache == nil {
		space.initializeCache = make(map[string]any)
	}
	if _, b := space.initializeCache[key]; b {
		return "", errors.New("'" + key + "' id already exists, repeated registration failed")
	}
	return key, nil
}

// DepKey 通过字段结构获取依赖 key
// 优先获取 tag属性的ref值
// 没有ref属性值 则获取包路径加类型全面
func DepKey(filed reflect.StructField) (string, bool) {
	depKey := ""
	// check 标识 通过 tag 方式初始化的 变量，必须要校验
	check := true
	switch filed.Type.Kind() {
	case reflect.Interface:
		// 如果字段是 接口类型 我们读取 ref 属性，ref 属性代表这个接口需要一个什么样的实现体，ref 也必须是 容器中可寻找到的属性
		if r, b := filed.Tag.Lookup("ref"); b && r != "" {
			depKey = r
		} else {
			// 没有 ref tag 则用字段名进行匹配
			depKey = filed.Name
			check = false
		}
	default:
		if r, b := filed.Tag.Lookup("ref"); b && r != "" {
			depKey = r
		} else {
			if filed.Type.Kind() == reflect.Ptr {
				depKey = fmt.Sprintf("%s-%s", filed.Type.Elem().PkgPath(), filed.Type.String())
			} else {
				depKey = fmt.Sprintf("%s-%s", filed.Type.PkgPath(), filed.Type.String())
			}
			check = false
		}
	}
	return depKey, check
}

// TypeKey 获取任意类型的 key
func TypeKey(t any) string {
	typeOf := reflect.TypeOf(t)
	baseType := ""
	if typeOf.Kind() == reflect.Ptr {
		baseType = fmt.Sprintf("%s-%s", typeOf.Elem().PkgPath(), typeOf.String())
	} else {
		baseType = fmt.Sprintf("%s-%s", typeOf.PkgPath(), typeOf.String())
	}
	return baseType
}

// Injection 反射初始化
// 把容器中对应的 value 赋值给 目标结构体的 field 字段
func Injection(field, value reflect.Value) error {
	if field.Type() == nil || value.Type() == nil {
		return nil
	}
	//校验是否可以分配赋值 类型方面的校验
	to := value.Type().AssignableTo(field.Type())
	if !to {
		//无法赋值 类型不匹配 结束运行
		return errors.New("the type is not assigned and cannot be assigned")
	}
	switch field.Kind() {
	case reflect.Interface:
		//如果字段是接口，我们 需要判断 value 是否实现了 filed 字段接口
		if value.Type().Implements(field.Type()) && value.Type().AssignableTo(field.Type()) {
			field.Set(value)
			return nil
		}
		return errors.New(value.Type().String() + "can not assignable to " + field.Type().String())
	case reflect.Pointer:
		if field.IsNil() {
			// 当前指针为空 设置指针指向value的地址
			if value.Elem().CanAddr() && field.CanSet() {
				//权限上面的校验
				field.Set(value.Elem().Addr())
			}
			return nil
		}
	}
	return nil
}
