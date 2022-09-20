package utils

import (
	"github.com/pkg/errors"
	"reflect"
	"time"
)

/*
处理反射中的复杂数据解析逻辑
func(value reflect.Value, data any) error 函数的 value 在结构体中视为一个基础数据处理，通过函数处理完成，代表结构体内对应的类型赋值完成
*/
type DataType func(value reflect.Value, data any) error

var dataType = map[string]DataType{
	reflect.TypeOf(time.Now()).String(): timeType,
}

// timeType 完成对时间 time.Time 的赋值操作
func timeType(value reflect.Value, data any) error {
	if s, ok := data.(string); !ok {
		return errors.New("Time.Time property initialization failed, please check whether the corresponding value format is correct")
	} else {
		parse, err := time.Parse("2006-04-02 15:04:04", s)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(parse))
	}
	return nil
}
