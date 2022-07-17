package utils

import (
	"errors"
	"github.com/iancoleman/strcase"
	"reflect"
	"strconv"
)

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
		//如果参数是接口
		ElementValue := field.Elem()
		return Injection(ElementValue, value)
	case reflect.Ptr:
		if field.IsNil() {
			// 当前指针为空 设置指针指向value的地址
			if value.Elem().CanAddr() {
				//权限上面的校验
				if field.CanSet() {
					field.Set(value.Elem().Addr())
				}
			}
			return nil
		}
	case reflect.Struct:
		for i := 0; i < field.NumField(); i++ {
			//查看字段名，通过字段名称进行一对一赋值
			name := field.Type().Field(i).Name
			f := field.FieldByName(name)
			_, b := value.Type().FieldByName(name)
			if !b {
				return errors.New("could not find matching field")
			}
			v := value.FieldByName(name)
			err := Injection(f, v)
			if err != nil {
				return err
			}
		}
	default:
		//校验是否可以设置值 权限上面的校验
		if field.CanSet() {
			field.Set(value)
		}
	}
	return nil
}

// StarAssignment
// data 传入 value 对应的 map[string]interface{}
func StarAssignment(value reflect.Value, data interface{}) error {
	switch value.Kind() {
	case reflect.Slice, reflect.Map, reflect.Struct:
		return Assignment(value, data)
	case reflect.Ptr:
		//需要先分配一个对应类型反射的值，这个值调用 Elem 获取对应指向的值才不会为空
		d := reflect.New(value.Type().Elem())
		//传递 指针指向的值，进行参数注入
		err := Assignment(d.Elem(), data)
		if err != nil {
			return err
		}
		//把分配的反射指针设置给参数
		value.Set(d)
	default:
		return Assignment(value, data)
	}
	return nil
}

// Assignment 递归对单个反射结构体进行赋值
// Assignment 带注入的参数,入参必须是值类型，指针类型赋值需要传递指针所指向的value，带注入的参数为map类型时候 arguments发射的map对应的接收类型必须是 map[string]interface{}
// value对应注入的 k/v
func Assignment(arguments reflect.Value, value interface{}) error {
	if value == nil || arguments.Type() == nil {
		return nil
	}
	var FieldName map[string]string
	//获取反射的类型
	t := arguments.Type()
	if t.Kind() == reflect.Struct {
		FieldName = make(map[string]string)
		//如果参数是结构体，初始化一份 字段名的对应表,该对应表，适配一些命名规范 蛇形,下划线什么的
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			snake := strcase.ToSnake(f.Name)
			camel := strcase.ToCamel(f.Name)
			kebab := strcase.ToKebab(f.Name)
			lowerCamel := strcase.ToLowerCamel(f.Name)
			FieldName[snake] = f.Name
			FieldName[camel] = f.Name
			FieldName[kebab] = f.Name
			FieldName[lowerCamel] = f.Name
			//支持标签定义，以防出现 strcase 库不稳定的时候 通过json来自定义
			if s, b := f.Tag.Lookup("json"); b {
				if s != "" {
					FieldName[s] = f.Name
				}
			}
		}
	}
	switch value.(type) {
	//如果参数为结构体，value对应为map 类型 此处的 map[string]interface{} 应当专属正对于 post 一类
	case map[string]interface{}:
		values := value.(map[string]interface{})
		if t.Kind() == reflect.Map {
			//如果 此处的 类型是map 则在此处就可以直接给 map 进行赋值 直接把 values 赋值给 value
			// 为什么在此处判断 map类型 ，避免走入下面的结构体解析逻辑中,这个赋值解析是以json的方式去对应 结构体 map[string]interface{} 在这个地方和结构体会有歧义
			//判断 map存储的值类型 是interface{} 就直接赋值
			// 此处进行修复 支持具体结构体，但是key的值必须是个字符串普通 key,请勿通过结构体形式作为key接收
			// 结构体形式的key，需要后续看情况。
			makeMap := reflect.MakeMap(arguments.Type())
			if err := AssignmentMap(makeMap, values); err != nil {
				return err
			}
			arguments.Set(makeMap)
			return nil
		}
		// 如果对应参数是 interface{} 将map传给它
		if t.Kind() == reflect.Interface {
			arguments.Set(reflect.ValueOf(value))
			return nil
		}
		//走到此处这个 arguments 必然是个结构体
		for k, v := range values {
			var field reflect.Value
			//判断结构体类型，及其校验结构体字段名对应表是否初始化
			if (t.Kind() == reflect.Struct) && FieldName != nil {
				if name, b := FieldName[k]; b {
					field = arguments.FieldByName(name)
				} else {
					//此处如果没有找到对应的字段名称，说明传递的json参数无法解析注入到 参数列表中
					//直接给前端返回错误提示，或者输出日志,此处直接日志返回，执行到具体处理器上参数将为零值
					return nil
				}
			} else {
				//校验结构体失败此处要么返回错误要么 panic，或者 return 放弃这个字段的初始化，处理器将会接收到零值
				return nil
			}
			if field.Type() == nil || v == nil { //v == nil 防止下面的 switch 走到 default中的 case reflect.Ptr 造成栈溢出
				return nil
			}
			switch v.(type) {
			case map[string]interface{}:
				//处理结构体类型字段
				if field.Kind() == reflect.Ptr {
					//指针类型，必须先分配内存,field 的类型为某个结构体的指针，先要获取到该结构体指针类型，指针指向的具体类型，然后为其分配New，New得到的才是 想要的指针类型
					v2 := reflect.New(field.Type().Elem())
					//获取指针的值，初始化复制先要拿到指针指向的值才可操作
					elem := v2.Elem()
					//初始化赋值
					if err := Assignment(elem, v); err != nil {
						return err
					}
					field.Set(v2)
				} else {
					if err := Assignment(field, v); err != nil {
						return err
					}
				}
			default:
				//处理普通字段属性
				if err := Assignment(field, v); err != nil {
					return err
				}
			}
		}
	//处理单个变量,结构体在上面的case 中应该已经被处理了
	default:
		//此处单一字段的基本类型转化，还是不能直接断言赋值，可能存在 传递json中的数字类型就是字符串形式，json解析为map的时候
		switch arguments.Kind() {
		case reflect.Interface:
			// 添加对 接口参数的支持
			arguments.Set(reflect.ValueOf(value))
		case reflect.String:
			variable := value.(string)
			arguments.SetString(variable)
		case reflect.Int:
			var variable int64
			switch value.(type) {
			case string:
				atoi, err := strconv.Atoi(value.(string))
				if err != nil {
					return errors.New("The reflection target parameter is of type int. The parameter you gave cannot be converted to type int. Please check the json format of the passed parameter.error:" + err.Error())
				}
				variable = int64(atoi)
			case float64:
				variable = int64(value.(float64))
			case int:
				variable = int64(value.(int))
			}
			arguments.SetInt(variable)
		case reflect.Float64:
			var variable float64
			switch value.(type) {
			case string:
				float, err := strconv.ParseFloat(value.(string), 64)
				if err != nil {
					return errors.New("The reflection target parameter is of type float64. The parameter you gave cannot be converted to type float64. Please check the json format of the passed parameter.error:" + err.Error())
				}
				variable = float
			case float64:
				variable = value.(float64)
			}
			arguments.SetFloat(variable)
		case reflect.Bool:
			var variable bool
			switch value.(type) {
			case string:
				parseBool, err := strconv.ParseBool(value.(string))
				if err != nil {
					return errors.New("The reflection target parameter is of type bool. The parameter you gave cannot be converted to type bool. Please check the json format of the passed parameter.error:" + err.Error())
				}
				variable = parseBool
			case bool:
				variable = value.(bool)
			}
			arguments.SetBool(variable)
		case reflect.Map:
			//通过前面的 case 内层的  reflect.Map 无法进入到这里

			//反射确定 value类型
		case reflect.Ptr:
			typ := arguments.Type().Elem()
			v := reflect.New(typ)
			elem := v.Elem()
			if err := Assignment(elem, value); err != nil {
				return err
			}
			arguments.Set(v)
		case reflect.Slice:
			elem := arguments.Type().Elem()
			slice := reflect.MakeSlice(arguments.Type(), 0, 0)
			arr, b := value.([]interface{})
			if !b {
				return errors.New("The reflection target parameter is of type slice. The parameter you gave cannot be converted to type slice. Please check the json format of the passed parameter.")
			}
			switch elem.Kind() {
			case reflect.Int:
				for _, element := range arr {
					var v int
					switch element.(type) {
					case string:
						atoi, err := strconv.Atoi(value.(string))
						if err != nil {
							return errors.New("The reflection target parameter is of type int. The parameter you gave cannot be converted to type int. Please check the json format of the passed parameter.error:" + err.Error())
						}
						v = atoi
					case float64:
						e := element.(float64)
						v = int(e)
					}
					slice = reflect.Append(slice, reflect.ValueOf(v))
				}
				arguments.Set(slice)
			case reflect.Float64:
				for _, element := range arr {
					var v float64
					switch element.(type) {
					case string:
						float, err := strconv.ParseFloat(element.(string), 64)
						if err != nil {
							return errors.New("The reflection target parameter is of type float64. The parameter you gave cannot be converted to type float64. Please check the json format of the passed parameter.error:" + err.Error())
						}
						v = float
					case float64:
						v = element.(float64)
					}
					slice = reflect.Append(slice, reflect.ValueOf(v))
				}
				arguments.Set(slice)
			case reflect.String:
				for _, element := range arr {
					e := element.(string)
					slice = reflect.Append(slice, reflect.ValueOf(e))
				}
				arguments.Set(slice)
			case reflect.Bool:
				for _, element := range arr {
					var bl bool
					switch element.(type) {
					case string:
						parseBool, err := strconv.ParseBool(element.(string))
						if err != nil {
							return errors.New("The reflection target parameter is of type bool. The parameter you gave cannot be converted to type bool." +
								" Please check the json format of the passed parameter.error:" + err.Error())
						}
						bl = parseBool
					case bool:
						bl = element.(bool)
					}
					slice = reflect.Append(slice, reflect.ValueOf(bl))
				}
				arguments.Set(slice)
			case reflect.Struct:
				for _, element := range arr {
					v := reflect.New(elem)
					v2 := v.Elem()
					if err := Assignment(v2, element); err != nil {
						return err
					}
					slice = reflect.Append(slice, v2)
				}
				arguments.Set(slice)
			case reflect.Ptr:
				for _, element := range arr {
					v := reflect.New(elem.Elem())
					v2 := v.Elem()
					if err := Assignment(v2, element); err != nil {
						return err
					}
					slice = reflect.Append(slice, v)
				}
				arguments.Set(slice)
			case reflect.Interface:
				for _, element := range arr {
					slice = reflect.Append(slice, reflect.ValueOf(element))
				}
				arguments.Set(slice)
			}
		}
	}
	return nil
}

// AssignmentMap 专门针对 map数据类型 进行解析
func AssignmentMap(arguments reflect.Value, value map[string]interface{}) error {
	if value == nil || arguments.Type() == nil {
		return nil
	}
	t := arguments.Type()
	//makeMap := reflect.MakeMap(t)
	//t.Elem() 获取map 存储的value类型
	switch t.Elem().Kind() {
	//检测 map存储的具体类型
	case reflect.Interface:
		//检测操map 存储value类型为接口，json解码的map刚好对应，所以可以直接通过反射赋值
		if t.Elem().Kind().String() == "interface" {
			for k, v := range value {
				key := reflect.New(t.Key()).Elem()
				if err := Assignment(key, k); err != nil {
					return err
				}
				val := reflect.New(t.Elem()).Elem()
				if err := Assignment(val, v); err != nil {
					return err
				}
				arguments.SetMapIndex(key, val)
			}
			return nil
		}
	case reflect.Ptr:
		//map 存储的数据为指针,通过反射获取到的类型也是指针，在分配内存的时候不能分配指针类型，要分配指针指向的类型才是正确结果
		for k, v := range value {
			key := reflect.New(t.Key()).Elem()
			if err := Assignment(key, k); err != nil {
				return err
			}
			v2 := reflect.New(t.Elem().Elem())
			val := v2.Elem()
			if err := Assignment(val, v); err != nil {
				return err
			}
			arguments.SetMapIndex(key, v2)
		}
		return nil
	case reflect.Int, reflect.Float64, reflect.String, reflect.Bool, reflect.Slice, reflect.Struct:
		for k, v := range value {
			key := reflect.New(t.Key()).Elem()
			if err := Assignment(key, k); err != nil {
				return err
			}
			val := reflect.New(t.Elem()).Elem()
			if err := Assignment(val, v); err != nil {
				return err
			}
			arguments.SetMapIndex(key, val)
		}
		return nil
	}
	return nil
}
