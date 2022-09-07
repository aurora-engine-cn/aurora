package sqlbuild

func (sql *SQL) Create(table string) *SQL {
	if sql.f {
		return sql
	}
	if !sql.c {
		sql.c = true
		sql.f = true
	}

	return sql
}
