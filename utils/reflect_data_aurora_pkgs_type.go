package utils

import (
	"errors"
	"fmt"
	"reflect"
)

//  处理 pkgs 包下面的自定义数据类型

func AuroraQueueType(value reflect.Value, data any) error {

	Queue := reflect.New(value.Type()).Elem()
	fmt.Println(Queue.Type().String())
	fmt.Println(Queue.NumMethod())

	EnQueueT, ok := Queue.Type().MethodByName("EnQueue")
	if !ok {
		return errors.New("not found EnQueue Method")
	}
	intype := EnQueueT.Type.In(0)
	arr, ok := data.([]interface{})
	if !ok {
		return errors.New("filed")
	}
	EnQueue := Queue.MethodByName("EnQueue")
	for _, v := range arr {
		elem := reflect.New(intype).Elem()
		err := Assignment(elem, v)
		if err != nil {
			return err
		}
		EnQueue.Call([]reflect.Value{elem})
	}
	value.Set(Queue)
	return nil
}

func AuroraQueuePointerType(value reflect.Value, data any) error {
	Queue := reflect.New(value.Type().Elem())
	EnQueue := Queue.MethodByName("EnQueue")
	intype := EnQueue.Type().In(0)
	arr, ok := data.([]interface{})
	if !ok {
		return errors.New("filed")
	}

	for _, v := range arr {
		elem := reflect.New(intype).Elem()
		err := Assignment(elem, v)
		if err != nil {
			return err
		}
		EnQueue.Call([]reflect.Value{elem})
	}
	value.Set(Queue)
	return nil
}

func AuroraStackType[T any](value reflect.Value, data any) {

}

func AuroraStackPointerType(value reflect.Value, data any) error {

	return nil
}
