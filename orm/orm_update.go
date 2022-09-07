package orm

import (
	"orm/sqlbuild"
	"reflect"
)

// Update 更新 单表
func (mapping *Mapping[T]) Update(condition, values T) int64 {
	defer mapping.unZero()
	updateSql := mapping.buildStructUpdateSql(condition, values)
	exec, err := mapping.DB.Exec(updateSql)
	if err != nil {
		panic(err.Error())
	}
	affected, err := exec.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	return affected
}

func (mapping *Mapping[T]) buildStructUpdateSql(condition, values T) string {
	var table string
	sql := sqlbuild.Sql()
	valueOf := reflect.ValueOf(condition)
	if valueOf.Kind() == reflect.Struct {
		valueOf = valueOf.Addr()
	}
	call := valueOf.MethodByName("Table").Call(nil)
	table = call[0].Interface().(string)
	sql.Update(table)
	set := mapping.buildStructUpdateSet(values)
	sql.Set(set)
	mapping.buildStructWhereSql(sql, condition)
	return sql.String()
}

func (mapping *Mapping[T]) buildStructUpdateSet(values T) map[string]any {
	valueOf := reflect.ValueOf(values)
	if valueOf.Kind() == reflect.Ptr && valueOf.Elem().Kind() == reflect.Struct {
		valueOf = valueOf.Elem()
	}
	set := make(map[string]any)
	var structField reflect.StructField
	var field reflect.Value
	var name string
	var fieldValue any
	for i := 0; i < valueOf.NumField(); i++ {
		// 待区分值问题
		field = valueOf.Field(i)
		if !mapping.zero && field.IsZero() {
			continue
		}
		fieldValue = field.Interface()
		structField = valueOf.Type().Field(i)
		name = mapping.analysisStructMapping(structField)
		//name = structField.Name
		//if get := structField.Tag.Get("column"); get != "" {
		//	name = get
		//}
		set[name] = fieldValue
	}
	return set
}

// UpdateMap map 形式构建 sql
func (mapping *Mapping[T]) UpdateMap(condition, values map[string]any) int64 {
	updateSql := mapping.buildMapUpdateSetSql(condition, values)
	exec, err := mapping.DB.Exec(updateSql)
	if err != nil {
		panic(err.Error())
	}
	affected, err := exec.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	return affected
}

func (mapping *Mapping[T]) buildMapUpdateSetSql(condition, values map[string]any) string {
	var table string
	sql := sqlbuild.Sql()
	valueOf := reflect.ValueOf(mapping.table)
	if valueOf.Kind() == reflect.Struct {
		valueOf = valueOf.Addr()
	}
	call := valueOf.MethodByName("Table").Call(nil)
	table = call[0].Interface().(string)
	sql.Update(table)
	sql.Set(values)
	mapping.buildMapWhereSql(sql, condition)
	return sql.String()
}
