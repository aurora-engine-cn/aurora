package sqlbuild

import "strings"

// Insert 创建插入语句
func (sql *SQL) Insert(table string, column ...string) *SQL {
	if sql.f {
		return sql
	}
	if !sql.i {
		sql.f = true
		sql.i = true
	}
	if sql.tables == nil {
		sql.tables = make([]string, 0)
	}
	if sql.columns == nil {
		sql.columns = make([]string, 0)
	}
	sql.tables = append(sql.tables, table)
	sql.columns = append(sql.columns, column...)
	sql.build()
	return sql
}

// Value 创建插入语句 value
func (sql *SQL) Value(value ...string) *SQL {
	if sql.values == nil {
		sql.values = make([]string, 0)
	}
	if join := strings.Join(value, ","); join != "" {
		join = "(" + join + ")"
		sql.values = append(sql.values, join)
	}
	sql.build()
	return sql
}
