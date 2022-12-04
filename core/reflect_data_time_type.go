package core

import (
	"errors"
	"gitee.com/aurora-engine/aurora/utils/strutils"
	"reflect"
	"time"
)

// TimeType 完成对时间 time.Time 的赋值操作
func TimeType(value reflect.Value, data any) error {
	timeV := ""
	switch data.(type) {
	case string:
		timeV = data.(string)
	// 适配get请求传递参数 直接解析入参的方式
	case map[string]interface{}:
		for _, timeValue := range data.(map[string]interface{}) {
			switch timeValue.(type) {
			case string:
				timeV = timeValue.(string)
			}
		}
	}
	if strutils.IsEmpty(timeV) {
		return errors.New("Time.Time property initialization failed, please check whether the corresponding value format is correct")
	}
	parse, err := time.Parse("2006-04-02 15:04:04", timeV)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(parse))
	return nil
}

// TimePointerType 完成对时间指针 *time.Time 的赋值操作
func TimePointerType(value reflect.Value, data any) error {
	timeV := ""
	switch data.(type) {
	case string:
		timeV = data.(string)
	// 适配get请求传递参数 直接解析入参的方式
	case map[string]interface{}:
		for _, timeValue := range data.(map[string]interface{}) {
			switch timeValue.(type) {
			case string:
				timeV = timeValue.(string)
			}
		}
	}
	if strutils.IsEmpty(timeV) {
		return errors.New("*Time.Time property initialization failed, please check whether the corresponding value format is correct")
	}

	parse, err := time.Parse("2006-04-02 15:04:04", timeV)
	if err != nil {
		return err
	}
	of := reflect.ValueOf(parse)
	//在次分配内存的原因在于 初始化的参数阶段虽然对整个结构体进行了分配，分配好的属性却是零值，对于指针的零值则需要额外的创建
	v := reflect.New(reflect.TypeOf(time.Time{}))
	v.Elem().Set(of)
	value.Set(v)

	return nil
}
