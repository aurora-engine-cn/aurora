package sqlbuild

import (
	"bytes"
	"strings"
)

// 构建 sql
func (sql *SQL) build() {
	sql.buffer = &bytes.Buffer{}
	if sql.s {
		// select构建
		// 构建 查询字段
		sql.columnsBuild()

		// 构建 普通表连接
		sql.tablesBuild()

		// 构建 查询条件
		sql.wheresBuild()

		// 构建 分组字段
		sql.groupsBuild()
		// 构建 分组条件
		sql.havingBuild()

		// 构建 排序条件
		sql.ordersBuild()

		// 分页条件构造
		sql.limitBuild()
		sql.offsetBuild()
	}

	if sql.i {
		// insert 构建
		sql.insertBuild()
		sql.valueBuild()
	}

	// update 构建
	if sql.u {
		sql.updateBuild()
		sql.setBuild()
		// 构建 查询条件
		sql.wheresBuild()
	}

	// delete 构建
	if sql.d {
		sql.deleteBuild()
		sql.wheresBuild()
	}
	// 结尾符
	sql.end()
}

func (sql *SQL) columnsBuild() {
	if join := strings.Join(sql.columns, ","); join != "" {
		sql.buffer.WriteString(SELECT)
		sql.buffer.WriteString(join)
	}
}

func (sql *SQL) tablesBuild() {
	if join := strings.Join(sql.tables, ","); join != "" {
		sql.buffer.WriteString(FROM)
		sql.buffer.WriteString(join)
	}
}

func (sql *SQL) wheresBuild() {
	var temp []string
	if sql.wheres != nil {
		temp = sql.wheres
		f := temp[len(temp)-1]
		if f == AND || f == OR {
			temp = temp[:len(sql.wheres)-1]
		}
	}
	if join := strings.Join(temp, " "); join != "" {
		sql.buffer.WriteString(WHERE)
		sql.buffer.WriteString(join)
	}
}

func (sql *SQL) groupsBuild() {
	if join := strings.Join(sql.groups, ","); join != "" {
		sql.buffer.WriteString(GROUP)
		sql.buffer.WriteString(join)
	}
}

func (sql *SQL) havingBuild() {
	if join := strings.Join(sql.having, AND); join != "" {
		if len(sql.having) > 1 {
			join = "(" + join + ")"
		}
		sql.buffer.WriteString(HAVING)
		sql.buffer.WriteString(join)
	}
}

func (sql *SQL) ordersBuild() {
	if join := strings.Join(sql.orders, ","); join != "" {
		sql.buffer.WriteString(ORDER)
		sql.buffer.WriteString(join)
	}
}

func (sql *SQL) limitBuild() {
	if join := strings.Join(sql.limit, ","); join != "" {
		sql.buffer.WriteString(LIMIT)
		sql.buffer.WriteString(join)
	}
}

func (sql *SQL) offsetBuild() {
	if join := strings.Join(sql.offset, ","); join != "" {
		sql.buffer.WriteString(OFFSET)
		sql.buffer.WriteString(join)
	}
}

// and 条件拼接
func (sql *SQL) and() {
	if sql.wheres == nil {
		sql.wheres = make([]string, 0)
	}
	sql.wheres = append(sql.wheres, AND)
}

// or 条件拼接
func (sql *SQL) or() {
	if sql.wheres == nil {
		sql.wheres = make([]string, 0)
	}
	sql.wheres = append(sql.wheres, OR)
}

// 构建插入语句 表名和字段
func (sql *SQL) insertBuild() {
	if join := strings.Join(sql.tables, ""); join != "" {
		sql.buffer.WriteString(INSERT)
		sql.buffer.WriteString(join)
		join = strings.Join(sql.columns, ",")
		join = "(" + join + ")"
		sql.buffer.WriteString(join)
	}
}

// 构建插入语句的 value部分
func (sql *SQL) valueBuild() {
	if join := strings.Join(sql.values, ","); join != "" {
		sql.buffer.WriteString(VALUES)
		sql.buffer.WriteString(join)
	}
}

func (sql *SQL) updateBuild() {
	if join := strings.Join(sql.tables, ","); join != "" {
		sql.buffer.WriteString(UPDATE)
		sql.buffer.WriteString(join)
	}
}

func (sql *SQL) setBuild() {
	if join := strings.Join(sql.set, ","); join != "" {
		sql.buffer.WriteString(SET)
		sql.buffer.WriteString(join)
	}
}

func (sql *SQL) deleteBuild() {
	if join := strings.Join(sql.tables, ","); join != "" {
		sql.buffer.WriteString(DELETE)
		sql.buffer.WriteString(FROM)
		sql.buffer.WriteString(join)
	}
}

func (sql *SQL) buildJoin() {

}

func (sql *SQL) buildLeft() {

}

func (sql *SQL) buildRight() {

}

// end 结束符
func (sql *SQL) end() {
	sql.buffer.WriteString(";")
}
