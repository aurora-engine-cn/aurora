package sqlbuild

import "strings"

// Where 创建查询语句，删除语句，更新语句的查询条件设置
func (sql *SQL) Where(where ...string) *SQL {
	if sql.wheres == nil {
		sql.wheres = make([]string, 0)
	}
	where = whereReplaceChart(where...)
	if join := strings.Join(where, AND); join != "" {
		if len(where) > 1 {
			join = "(" + join + ")"
		}
		sql.wheres = append(sql.wheres, join)
		sql.and()
	}
	sql.build()
	return sql
}

// And 与条件连接符
func (sql *SQL) And() *SQL {
	// 删除 最后一个元素 添加一个新的 连接符
	sql.wheres = append(sql.wheres[:len(sql.wheres)-1], " AND ")
	sql.build()
	return sql
}

// Or 或条件连接符
func (sql *SQL) Or() *SQL {
	sql.wheres = append(sql.wheres[:len(sql.wheres)-1], " OR ")
	sql.build()
	return sql
}
