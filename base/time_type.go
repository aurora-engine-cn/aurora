package base

import (
	"errors"
	"reflect"
	"time"
)

// TimeType 完成对时间 time.Time 的赋值操作
func TimeType(value reflect.Value, data any) error {
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

// TimeTypePointer 完成对时间指针 *time.Time 的赋值操作
func TimeTypePointer(value reflect.Value, data any) error {
	if s, ok := data.(string); !ok {
		return errors.New("*Time.Time property initialization failed, please check whether the corresponding value format is correct")
	} else {
		parse, err := time.Parse("2006-04-02 15:04:04", s)
		if err != nil {
			return err
		}
		of := reflect.ValueOf(parse)
		//在次分配内存的原因在于 初始化的参数阶段虽然对整个结构体进行了分配，分配好的属性却是零值，对于指针的零值则需要额外的创建
		v := reflect.New(reflect.TypeOf(time.Time{}))
		v.Elem().Set(of)
		value.Set(v)
	}
	return nil
}
