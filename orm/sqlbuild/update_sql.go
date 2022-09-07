package sqlbuild

import (
	"fmt"
	"github.com/druidcaesa/ztool"
	"time"
)

// Update 创建更新语句
func (sql *SQL) Update(table string) *SQL {
	if sql.f {
		return sql
	}
	if !sql.u {
		sql.u = true
		sql.f = true
	}
	if sql.tables == nil {
		sql.tables = make([]string, 0)
	}
	sql.tables = append(sql.tables, table)
	sql.build()
	return sql
}

// Set 设置更新的 字段
func (sql *SQL) Set(column map[string]any) *SQL {
	if sql.set == nil {
		sql.set = make([]string, 0)
	}
	for k, v := range column {
		var value string
		switch v.(type) {
		case string:
			value = fmt.Sprintf("%s='%s'", k, v)
		case int:
			value = fmt.Sprintf("%s=%d", k, v)
		case float64:
			value = fmt.Sprintf("%s=%f", k, v)
		case time.Time:
			format := ztool.DateUtils.SetTime(v.(time.Time)).Format()
			value = fmt.Sprintf("%s='%s'", k, format)
		}
		value = replaceSet(value)
		sql.set = append(sql.set, value)
	}
	sql.build()
	return sql
}
