package sqlbuild

// Select 查询字段设置
func (sql *SQL) Select(column ...string) *SQL {
	if sql.f {
		return sql
	}
	if !sql.s {
		sql.f = true
		sql.s = true
	}
	if sql.columns == nil {
		sql.columns = make([]string, 0)
	}
	column = selectReplaceChart(column...)
	sql.columns = append(sql.columns, column...)

	sql.build()
	return sql
}

// From 查询表设置
// 一般设置单个表
// 多个表拼接为笛卡尔积形式(全连接)
func (sql *SQL) From(table ...string) *SQL {
	if !sql.s {
		return sql
	}
	if sql.tables == nil {
		sql.tables = make([]string, 0)
	}
	sql.tables = append(sql.tables, table...)
	sql.build()
	return sql
}

// Group 分组字段设置
func (sql *SQL) Group(column ...string) *SQL {
	sql.groups = make([]string, 0)
	sql.groups = append(sql.groups, column...)
	sql.build()
	return sql
}

// Having 分组条件字段
// 分组条件仅支持 一个条件
func (sql *SQL) Having(having ...string) *SQL {
	if sql.having == nil {
		sql.having = make([]string, 0)
	}
	sql.having = append(sql.having, having...)
	sql.build()
	return sql
}

// Order 排序字段
// 需要手动添加排序关键字
func (sql *SQL) Order(order ...string) *SQL {
	if sql.orders == nil {
		sql.orders = make([]string, 0)
	}
	sql.orders = append(sql.orders, order...)
	sql.build()
	return sql
}

// Limit 分页条件设置
// 一个参数为设置条数 limit a
// 两个参数则设置 limit a,b
func (sql *SQL) Limit(limit ...string) *SQL {
	if sql.limit == nil {
		sql.limit = make([]string, 0)
	}
	sql.limit = append(sql.limit, limit...)
	sql.build()
	return sql
}

// Offset 分页设置结束位置
// 传递一个参数即可
func (sql *SQL) Offset(offset ...string) *SQL {
	if sql.offset == nil {
		sql.offset = make([]string, 0)
	}
	sql.offset = append(sql.offset, offset...)
	sql.build()
	return sql
}
