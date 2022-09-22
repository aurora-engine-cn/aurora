package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

//处理反射中的复杂数据解析逻辑，提供这样处理数据的初衷是因为 go 中的某些类型，比如说时间 内部是存在不可导出的字段，导致了递归解析中会出现问题
//这一类数据我们一般是当作一个基础字段来使用，并不关注内部的赋值过程，只需把对应的值按照正确的逻辑赋值给变量即可，至于赋值的逻辑就可以通过 DataType 来实现
//在全局变量 BaseType 中添加一个k/v 来对想要支持的字段进行解析操作。
//这个数据的解析操作只支持在结构体内的字段属性，作为web接口函数，支持的基本数据类型已经满足绝大部分使用，并且添加入参级别的基础类型解析，需要考虑值类型和指针类型的适配相对繁琐，在此整个请求解析中只支持结构体级别的自定义解析

// 存在的问题
// 1. 当前的处理方式仅能满足不是泛型方式的类型，因为通过（类型字符串-->reflect.TypeOf(t).String() ）的方式对同一个类型的不同泛型比如( Type[int],Type[any]...) 都是以一个新的类型字符串标识 无法对多种形式的泛型做处理。
// 	  不能通过类型字符串方式指定泛型的解析逻辑 但是 按照反射类型的方式可以正常处理可导出的泛型类型其根本原因是递归解析中不依赖类型字符串作为处理逻辑(对应第三方数据结构库会存在保护变量不导出的情况，将成为一个解析阻碍)
// 2. 对泛型的参数支持只能支持到 别名的方式声明的泛型定义 比如 type ArrayList[T any] []T  type Map[K comparable, V any] map[K]V  此类的泛型 还是保留者基础数据类型的特性.
// 3. 当前 base 包仅处理不是泛型的自定义类型解析

func init() {
	BaseType = map[string]DataType{
		// 加载时间类型基础变量
		TypeKey(time.Time{}): TimeType,

		// 加载时间指针类型基础变量
		TypeKey(&time.Time{}): TimePointerType,

		// pkgs 队列类型参数支持
		"gitee.com/aurora-engine/aurora/pkgs/queue-*queue.Queue": AuroraQueuePointerType,
		// pkgs 栈类型参数支持
		"gitee.com/aurora-engine/aurora/pkgs/stack-*stack.Stack": AuroraStackPointerType,
	}
}

// DataType 函数定义反射赋值逻辑
// value : 是在一个结构体内的字段反射，通过该函数可以对这个字段进行初始化赋值
// data  : 是value对应的具体参数值，可能是字符串，切片，map
type DataType func(value reflect.Value, data any) error

// BaseType 存储了请求参数解析过程中对结构体内部字段类型的自定义支持，添加到 Type 中的类型在 控制器参数校验时候会自动跳过
// key : 通过对类型的反射取到的类型名称
// value : 定义了对应该类型的解析逻辑
var BaseType map[string]DataType

// TypeKey 通过反射得到一个类型的类型字符串
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

func BaseTypeKey(v reflect.Value) string {
	baseType := ""
	if v.Kind() == reflect.Ptr {
		baseType = fmt.Sprintf("%s-%s", v.Type().Elem().PkgPath(), v.Type().String())
	} else {
		baseType = fmt.Sprintf("%s-%s", v.Type().PkgPath(), v.Type().String())
	}
	if index := strings.Index(baseType, "["); index != -1 {
		baseType = baseType[:index]
	}
	return baseType
}
