package sqlbuild

import (
	"bytes"
)

func Sql() *SQL {
	return &SQL{}
}

type SQL struct {
	columns, tables, wheres, groups, having, orders, offset, limit []string
	set                                                            []string
	insert, values                                                 []string
	joins, lefts, rights                                           []string
	// 构建标识
	i, s, d, u, c bool

	// 表连接标识
	left, right, join bool
	// 标识 sql 类型
	f bool
	// sql 语句
	buffer *bytes.Buffer
}

func (sql *SQL) String() string {
	return sql.buffer.String()
}
