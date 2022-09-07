package orm

import (
	"fmt"
	"github.com/druidcaesa/ztool"
	"orm/sqlbuild"
	"reflect"
	"strings"
	"time"
)

// 更新 删除 where 构建 api

func (mapping *Mapping[T]) buildStructWhereSql(sql *sqlbuild.SQL, condition T) {
	valueOf := reflect.ValueOf(condition)
	if valueOf.Kind() == reflect.Ptr && valueOf.Elem().Kind() == reflect.Struct {
		valueOf = valueOf.Elem()
	}
	// 字段信息
	var structField reflect.StructField
	// 字段值信息
	var field reflect.Value
	// 条件构造
	var where, column string
	// 字段具体值
	var fieldV any

	for i := 0; i < valueOf.NumField(); i++ {
		field = valueOf.Field(i)
		structField = valueOf.Type().Field(i)
		if !mapping.zero && field.IsZero() {
			continue
		}
		// 针对一级指针参数做处理
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		column = mapping.analysisStructMapping(structField)
		fieldV = field.Interface()
		switch fieldV.(type) {
		case string:
			where = fmt.Sprintf("%s = '%s'", column, fieldV)
		case int:
			where = fmt.Sprintf("%s = %d", column, fieldV)
		case float64:
			where = fmt.Sprintf("%s = %f", column, fieldV)
		case time.Time:
			where = fmt.Sprintf("%s = '%s'", column, ztool.DateUtils.SetTime(fieldV.(time.Time)).Format())
		case *time.Time:
			where = fmt.Sprintf("%s = '%s'", column, ztool.DateUtils.SetTime(*fieldV.(*time.Time)).Format())
		}
		sql.Where(where)
	}
}

// 解析结构体和数据库上面的映射信息
func (mapping *Mapping[T]) analysisStructMapping(structField reflect.StructField) string {
	var column string
	// 默认采用字段名
	column = structField.Name

	// 解析支持的 column tag
	if get := structField.Tag.Get("column"); get != "" {
		attr := columnFormat(get)
		column = attr[0]
	}
	// 其他支持项...

	// 支持 gorm column 属性
	if get := structField.Tag.Get("gorm"); get != "" {
		split := strings.Split(get, ";")
		for _, v := range split {
			if strings.Contains(v, "column") {
				c := strings.Split(v, ":")
				return c[1]
			}
		}
	}
	return column
}
func columnFormat(definition string) []string {
	split := strings.Split(definition, " ")
	attr := make([]string, 0)
	for i := 0; i < len(split); i++ {
		if split[i] != "" {
			attr = append(attr, split[i])
		}
	}
	return attr
}

// 构建 map k/v 条件
func (mapping *Mapping[T]) buildMapWhereSql(sql *sqlbuild.SQL, condition map[string]any) {
	var where string
	for k, v := range condition {
		switch v.(type) {
		case string:
			where = fmt.Sprintf("%s = '%s'", k, v)
		case int:
			where = fmt.Sprintf("%s = %d", k, v)
		case float64:
			where = fmt.Sprintf("%s = %f", k, v)
		case time.Time:
			where = fmt.Sprintf("%s = '%s'", k, v)
		case *time.Time:
			where = fmt.Sprintf("%s = '%s'", k, v)

		}
		sql.Where(where)
	}
}
