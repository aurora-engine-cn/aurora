package orm

import (
	"database/sql"
	"fmt"
	"gitee.com/aurora-engine/aurora/orm/sqlbuild"
	"github.com/druidcaesa/ztool"
	"log"
	"reflect"
	"time"
)

// Select 根据结构体查询单条记录
// 默认不对零值进行构建，结构体没有 column 字段 将默认采用结构体字段名作为列名构建
// 查询条件如果会在多条也只会返回其中一条为结果，如需擦汗寻多条请看 Selects 和 SelectMaps
func (mapping *Mapping[T]) Select(condition T) T {
	defer mapping.unZero()
	var v T
	selectSql := mapping.buildStructSelectSql(condition)
	row, err := mapping.DB.Query(selectSql)
	if err != nil {
		log.Println(err.Error())
		return v
	}
	slice := mapping.resultSlice(row, condition)
	if slice != nil && len(slice) > 0 {
		v = slice[0]
	}
	return v
}

func (mapping *Mapping[T]) Selects(condition T) []T {
	defer mapping.unZero()
	// 准备返回值
	var v []T
	// 生成 Select Sql 语句
	selectSql := mapping.buildStructSelectSql(condition)
	fmt.Println(selectSql)
	// 执行查询
	row, err := mapping.DB.Query(selectSql)
	if err != nil {
		log.Println(err.Error())
		return v
	}
	// 解析结果集
	v = mapping.resultSlice(row, condition)
	return v
}

// SelectMap k/v 条件查询单条记录
func (mapping *Mapping[T]) SelectMap(condition map[string]any) T {
	var v T
	selectSql := mapping.buildMapSelectSql(condition)
	row, err := mapping.DB.Query(selectSql)
	if err != nil {
		log.Println(err.Error())
		return v
	}
	slice := mapping.resultSlice(row, v)
	if slice != nil && len(slice) > 0 {
		v = slice[0]
	}
	return v
}

// SelectMaps k/v 条件查询多条记录
func (mapping *Mapping[T]) SelectMaps(condition map[string]any) []T {
	var v []T
	selectSql := mapping.buildMapSelectSql(condition)
	row, err := mapping.DB.Query(selectSql)
	if err != nil {
		log.Println(err.Error())
		return v
	}
	return mapping.resultSlice(row, mapping.table)
}

// 生成 结构体查询sql语句
func (mapping *Mapping[T]) buildStructSelectSql(condition T) string {
	var of reflect.Value
	var table string
	of = reflect.ValueOf(condition)
	if of.Kind() == reflect.Struct {
		of = reflect.ValueOf(&condition)
	}
	call := of.MethodByName("Table").Call(nil)
	table = call[0].Interface().(string)
	statement := sqlbuild.Sql()
	statement.Select("*")
	statement.From(table)
	mapping.buildStructWhereSql(statement, condition)
	return statement.String()
}

// 生成结构体 映射匹配
func structMapping(s any) map[string]string {
	mapp := make(map[string]string)
	var of reflect.Type
	if reflect.TypeOf(s).Kind() == reflect.Ptr {
		of = reflect.TypeOf(s).Elem()
	} else {
		of = reflect.TypeOf(s)
	}
	for i := 0; i < of.NumField(); i++ {
		field := of.Field(i)
		mapp[field.Name] = field.Name
		if get := field.Tag.Get("column"); get != "" {
			mapp[get] = field.Name
		}
	}
	return mapp
}

// 查询结果集封装
func (mapping *Mapping[T]) resultSlice(row *sql.Rows, condition T) []T {
	of := reflect.ValueOf(row)
	// 确定数据库 列顺序 排列扫描顺序
	columns, err := row.Columns()
	if err != nil {
		panic(err.Error())
	}
	// 解析结构体 映射字段
	// 拿到 scan 方法
	scan := of.MethodByName("Scan")
	next := of.MethodByName("Next")
	t := reflect.SliceOf(reflect.TypeOf(condition))
	result := reflect.MakeSlice(t, 0, 0)
	for (next.Call(nil))[0].Interface().(bool) {
		var value, unValue reflect.Value
		if reflect.TypeOf(condition).Kind() == reflect.Ptr {
			//创建一个 接收结果集的变量
			value = reflect.New(reflect.TypeOf(condition).Elem())
			unValue = value.Elem()
		} else {
			value = reflect.New(reflect.TypeOf(condition))
			value = value.Elem()
			unValue = value
		}
		// 创建 接收器
		values, fieldIndexMap := mapping.buildScan(unValue, columns)
		// 执行扫描, 执行结果扫描，不处理error 扫码结果类型不匹配，默认为零值
		scan.Call(values)
		// 迭代是否有特殊结构体 主要对 时间类型做了处理
		mapping.scanWrite(values, fieldIndexMap)
		// 添加结果集
		result = reflect.Append(result, value)
	}
	return result.Interface().([]T)
}

// 构建结构体接收器
func (mapping *Mapping[T]) buildScan(value reflect.Value, columns []string) ([]reflect.Value, map[int]reflect.Value) {
	// Scan 函数调用参数列表,接收器存储的都是指针类 反射的指针类型
	values := make([]reflect.Value, 0)
	// 存储的 也将是指针的反射形式
	fieldIndexMap := make(map[int]reflect.Value)
	// 创建 接收器
	for _, column := range columns {
		// 通过结构体映射找到 数据库映射到结构体的字段名
		name := mapping.columns[column]
		// 找到对应的字段
		byName := value.FieldByName(name)
		// 检查 接收参数 如果是特殊参数 比如结构体，时间类型的情况需要特殊处理 当前仅对时间进行特殊处理 ,获取当前 参数的 values 索引 并保存替换
		field := byName.Interface()
		switch field.(type) {
		case time.Time:
			// 记录特殊 值的索引 并且替换掉
			index := len(values)
			fieldIndexMap[index] = byName.Addr()
			// 替换 默认使用空字符串去接收
			values = append(values, reflect.New(reflect.TypeOf("")))
			continue
		case *time.Time:
			// 记录特殊 值的索引 并且替换掉
			index := len(values)
			fieldIndexMap[index] = byName.Addr()
			// 替换 使用空字符串去接收
			values = append(values, reflect.New(reflect.TypeOf("")))
			continue
		}
		values = append(values, byName.Addr())
	}
	return values, fieldIndexMap
}

// 对 buildScan 函数构建阶段存在特殊字段的处理 进行回写到指定的结构体位置
func (mapping *Mapping[T]) scanWrite(values []reflect.Value, fieldIndexMap map[int]reflect.Value) {
	// 迭代是否有特殊结构体 主要对 时间类型做了处理
	for k, v := range fieldIndexMap {
		// 拿到 特殊结构体对应的 值
		mapV := values[k]
		structField := v.Interface()
		switch structField.(type) {
		case *time.Time:
			// 吧把对应的 mappv 转化为 time.Time
			mappvalueString := mapV.Elem().Interface().(string)
			parse, err := ztool.DateUtils.Parse(mappvalueString)
			if err != nil {
				panic(err)
			}
			t2 := parse.Time()
			valueOf := reflect.ValueOf(t2)
			//设置该指针指向的值
			v.Elem().Set(valueOf)
		case **time.Time:
			// 吧把对应的 mappv 转化为 time.Time
			mappvalueString := mapV.Elem().Interface().(string)
			parse, err := ztool.DateUtils.Parse(mappvalueString)
			if err != nil {
				panic(err)
			}
			t2 := parse.Time()
			valueOf := reflect.ValueOf(&t2)
			//设置该指针指向的值
			v.Elem().Set(valueOf)
		}
	}
}

func (mapping *Mapping[T]) buildMapSelectSql(condition map[string]any) string {
	sql := sqlbuild.Sql()
	var of reflect.Value
	var table string
	of = reflect.ValueOf(mapping.table)
	if of.Kind() == reflect.Struct {
		of = reflect.ValueOf(&mapping.table)
	}
	call := of.MethodByName("Table").Call(nil)
	table = call[0].Interface().(string)
	sql.Select("*")
	sql.From(table)
	mapping.buildMapWhereSql(sql, condition)
	return sql.String()
}
