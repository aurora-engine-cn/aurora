package sqlbuild

// Delete 创建删除语句
func (sql *SQL) Delete(table string) *SQL {
	if sql.f {
		return sql
	}
	if !sql.d {
		sql.f = true
		sql.d = true
	}
	if sql.tables == nil {
		sql.tables = make([]string, 0)
	}
	sql.tables = append(sql.tables, table)
	sql.build()
	return sql
}
