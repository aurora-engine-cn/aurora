package base

import (
	"reflect"
	"time"
)

//处理反射中的复杂数据解析逻辑，提供这样处理数据的初衷是因为 go 中的某些类型，比如说时间 内部是存在不可导出的字段，导致了递归解析中会出现问题
//这一类数据我们一般是当作一个基础字段来使用，并不关注内部的赋值过程，只需把对应的值按照正确的逻辑赋值给变量即可，至于赋值的逻辑就可以通过 DataType 来实现
//在全局变量 BaseType 中添加一个k/v 来对想要支持的字段进行解析操作。
//这个数据的解析操作只支持在结构体内的字段属性，作为web接口函数，支持的基本数据类型已经满足绝大部分使用，并且添加入参级别的基础类型解析，需要考虑值类型和指针类型的适配相对繁琐，在此整个请求解析中只支持结构体级别的自定义解析

// DataType 函数定义反射赋值逻辑
// value : 是在一个结构体内的字段反射，通过该函数可以对这个字段进行初始化赋值
// data  : 是value对应的具体参数值
type DataType func(value reflect.Value, data any) error

// Type 存储了请求参数解析过程中对结构体内部字段类型的自定义支持
// key : 通过对类型的反射取到的类型名称
// value : 定义了对应该类型的解析逻辑
var Type = map[string]DataType{
	// 加载时间类型基础变量
	reflect.TypeOf(time.Time{}).String(): TimeType,

	// 加载时间指针类型基础变量
	reflect.TypeOf(&time.Time{}).String(): TimeTypePointer,
}
