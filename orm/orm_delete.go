package orm

import (
	"orm/sqlbuild"
	"reflect"
)

// Delete 删除单张表数据
func (mapping *Mapping[T]) Delete(delete T) int64 {
	defer mapping.unZero()
	deleteSql := mapping.buildStructDeleteSql(delete)
	exec, err := mapping.DB.Exec(deleteSql)
	if err != nil {
		panic(err.Error())
	}
	affected, err := exec.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	return affected
}

func (mapping *Mapping[T]) buildStructDeleteSql(delete T) string {
	sql := sqlbuild.Sql()
	var valueOf reflect.Value
	var table string
	valueOf = reflect.ValueOf(delete)
	if valueOf.Kind() == reflect.Struct {
		valueOf = valueOf.Addr()
	}
	call := valueOf.MethodByName("Table").Call(nil)
	table = call[0].Interface().(string)
	sql.Delete(table)
	mapping.buildStructWhereSql(sql, delete)
	return sql.String()
}

func (mapping *Mapping[T]) DeleteMap(delete map[string]any) int64 {
	defer mapping.unZero()
	deleteSql := mapping.buildMapDeleteSql(delete)
	exec, err := mapping.DB.Exec(deleteSql)
	if err != nil {
		panic(err.Error())
	}
	affected, err := exec.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	return affected
}

func (mapping *Mapping[T]) buildMapDeleteSql(delete map[string]any) string {
	sql := sqlbuild.Sql()
	var valueOf reflect.Value
	var table string
	valueOf = reflect.ValueOf(delete)
	if valueOf.Kind() == reflect.Struct {
		valueOf = valueOf.Addr()
	}
	call := valueOf.MethodByName("Table").Call(nil)
	table = call[0].Interface().(string)
	sql.Delete(table)
	mapping.buildMapWhereSql(sql, delete)
	return sql.String()
}
