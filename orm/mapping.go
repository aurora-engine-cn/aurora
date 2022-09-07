package orm

import (
	"database/sql"
	"reflect"
)

// CreateMapping 创建映射器
func CreateMapping[T any](db *sql.DB) *Mapping[T] {
	var t T
	if reflect.TypeOf(t).Kind() == reflect.Ptr {
		elem := reflect.New(reflect.TypeOf(t).Elem())
		t = elem.Interface().(T)
	} else {
		elem := reflect.New(reflect.TypeOf(t)).Elem()
		t = elem.Interface().(T)
	}
	m := &Mapping[T]{
		table:   t,
		columns: structMapping(t),
		DB:      db,
	}
	return m
}

type Mapping[T any] struct {
	zero    bool              // 是否启用零值 true 启用  false 不起用  默认不启用
	table   T                 //表信息
	columns map[string]string //字段映射关系
	DB      *sql.DB           //数据库连接
}

// Zero 启用零值
func (mapping *Mapping[T]) Zero() *Mapping[T] {
	mapping.zero = true
	return mapping
}

// 关闭启用
func (mapping *Mapping[T]) unZero() {
	mapping.zero = false
}
