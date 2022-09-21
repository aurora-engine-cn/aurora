package pkgstype

import (
	"reflect"
)

// base_aurora_pkgs_type.go 处理 pkgs 包下面的自定义数据类型

func AuroraQueueType(value reflect.Value, data any) error {

	return nil
}

func AuroraQueuePointerType[T any](value reflect.Value, data any) {

}

func AuroraStackType[T any](value reflect.Value, data any) {

}

func AuroraStackPointerType(value reflect.Value, data any) {

}
