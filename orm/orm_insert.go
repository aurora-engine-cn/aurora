package orm

import (
	"fmt"
	"gitee.com/aurora-engine/aurora/orm/sqlbuild"
	"github.com/druidcaesa/ztool"
	"reflect"
	"time"
)

/*
	插入数据库操作 对自增长主键暂不支持
*/

// Insert 插入一批数据
func (mapping *Mapping[T]) Insert(values ...T) int64 {
	defer mapping.unZero()
	sql := mapping.insertBuildStructSql(values...)
	exec, err := mapping.DB.Exec(sql)
	if err != nil {
		panic(err.Error())
	}
	affected, err := exec.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	return affected
}

func (mapping *Mapping[T]) insertBuildStructSql(values ...T) string {
	var table string
	sql := sqlbuild.Sql()
	column := mapping.insertColumn()
	valueOf := reflect.ValueOf(mapping.table)
	if valueOf.Kind() == reflect.Ptr && valueOf.Type().Elem().Kind() == reflect.Struct {
		tv := valueOf.MethodByName("Table").Call(nil)
		table = tv[0].Interface().(string)
	}
	if valueOf.Kind() == reflect.Struct {
		tv := reflect.ValueOf(&mapping.table).MethodByName("Table").Call(nil)
		table = tv[0].Interface().(string)
	}
	sql.Insert(table, column...)
	value := mapping.insertSetValue(values...)
	for i := 0; i < len(value); i++ {
		sql.Value(value[i]...)
	}
	return sql.String()
}

// 解析结构体 Insert 字段
func (mapping *Mapping[T]) insertColumn() []string {
	var structInfo reflect.Type
	typeOf := reflect.TypeOf(mapping.table)
	if typeOf.Kind() == reflect.Ptr {
		structInfo = typeOf.Elem()
	}
	if typeOf.Kind() == reflect.Struct {
		structInfo = typeOf
	}
	if structInfo == nil {
		return nil
	}
	columns := make([]string, 0)
	column := ""
	for i := 0; i < structInfo.NumField(); i++ {
		field := structInfo.Field(i)
		// 待添加对自增长主键校验 存在自增长主键 构建插入语句的时候不指定到语句中
		column = mapping.analysisStructMapping(field)
		columns = append(columns, column)
	}
	return columns
}

// 生成结构体插入value
func (mapping *Mapping[T]) insertSetValue(values ...T) [][]string {
	setValues := make([][]string, 0)
	for i := 0; i < len(values); i++ {
		set := make([]string, 0)
		ofValue := reflect.ValueOf(values[i])
		if ofValue.Kind() == reflect.Ptr && ofValue.Type().Elem().Kind() == reflect.Struct {
			ofValue = ofValue.Elem()
		}
		if ofValue.Kind() != reflect.Struct {
			return setValues
		}
		for j := 0; j < ofValue.NumField(); j++ {
			field := ofValue.Field(j)
			// 待添加对自增长主键校验 存在自增长主键 构建插入语句的时候不指定到语句中
			v := field.Interface()
			setv := ""
			switch v.(type) {
			case string:
				setv = fmt.Sprintf("'%s'", v)
			case int:
				setv = fmt.Sprintf("%d", v)
			case float64:
				setv = fmt.Sprintf("%f", v)
			case time.Time:
				setv = fmt.Sprintf("'%s'", ztool.DateUtils.SetTime(v.(time.Time)).Format())
			}
			set = append(set, setv)
		}
		setValues = append(setValues, set)
	}
	return setValues
}

func (mapping *Mapping[T]) InsertMap(values ...map[string]any) int64 {
	sql := mapping.insertBuildMapSql(values...)
	exec, err := mapping.DB.Exec(sql)
	if err != nil {
		panic(err.Error())
	}
	affected, err := exec.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	return affected
}

func (mapping *Mapping[T]) insertBuildMapSql(values ...map[string]any) string {
	var table string
	sql := sqlbuild.Sql()
	column := mapping.insertColumn()
	valueOf := reflect.ValueOf(mapping.table)
	if valueOf.Kind() == reflect.Ptr && valueOf.Type().Elem().Kind() == reflect.Struct {
		tv := valueOf.MethodByName("Table").Call(nil)
		table = tv[0].Interface().(string)
	}
	if valueOf.Kind() == reflect.Struct {
		tv := reflect.ValueOf(&mapping.table).MethodByName("Table").Call(nil)
		table = tv[0].Interface().(string)
	}
	column, value := mapping.insertMapSetColumnValue(values...)
	sql.Insert(table, column...)
	for i := 0; i < len(value); i++ {
		sql.Value(value[i]...)
	}
	return sql.String()
}

// 生成 map 插入语句的 列和值
func (mapping *Mapping[T]) insertMapSetColumnValue(values ...map[string]any) ([]string, [][]string) {
	// 获取 字段顺序列表
	columns := mapping.insertColumn()
	setValues := make([][]string, 0)
	// 存储待插入字段索引
	Column := make(map[int]struct{})
	for i := 0; i < len(values); i++ {
		set := make([]string, 0)
		setv := ""
		mv := values[i]
		for j := 0; j < len(columns); j++ {
			if v, b := mv[columns[j]]; b {
				Column[j] = struct{}{}
				switch v.(type) {
				case string:
					setv = fmt.Sprintf("'%s'", v)
				case int:
					setv = fmt.Sprintf("%d", v)
				case float64:
					setv = fmt.Sprintf("%f", v)
				case time.Time:
					setv = fmt.Sprintf("'%s'", ztool.DateUtils.SetTime(v.(time.Time)).Format())
				}
				set = append(set, setv)
			}
		}
		setValues = append(setValues, set)
	}
	// 生成插入 字段
	column := make([]string, 0)
	for index, _ := range Column {
		column = append(column, columns[index])
	}
	return column, setValues
}
