package build

import (
	"strings"
)

const (
	and = ") AND ("
	or  = ") OR ("
)

type Value interface{}

type SQL struct {
	column, table, join, left_join, right_join, on, where, group, order         []string
	condition, joins, left, right, o_n, selects, tables, wheres, groups, orders bool
}

func New() *SQL {
	return &SQL{
		column:     []string{},
		table:      []string{},
		join:       []string{},
		left_join:  []string{},
		right_join: []string{},
		where:      []string{},
	}
}

func (s *SQL) SELECT(column ...string) {
	s.selects = true
	s.column = append(s.column, column...)
}

func (s *SQL) FROM(table ...string) {
	s.tables = true
	s.table = append(s.table, table...)
}

func (s *SQL) LEFT_JOIN(join ...string) {
	s.left = true
	s.left_join = append(s.left_join, join...)
}

func (s *SQL) RIGHT_JOIN(join ...string) {
	s.right = true
	s.right_join = append(s.right_join, join...)
}

func (s *SQL) ON(on ...string) {
	s.o_n = true
	s.on = append(s.on, on...)
}

func (s *SQL) WHERE(sql string, value ...Value) {
	s.wheres = true
	analysis(sql, value...)
}

func (s *SQL) AND() {
	s.condition = true
	s.where = append(s.where, and)
}

func (s *SQL) OR() {
	s.condition = true
	s.where = append(s.where, or)
}

func (s *SQL) GROUP_BY(group ...string) {
	s.groups = true

}

func (s *SQL) WITH() {

}

func (s *SQL) HAVING() {

}

func (s *SQL) ORDER_BY() {

}

func (s *SQL) LIMIT(limit string) {

}

func (s *SQL) OFFSET(offset string) {

}

func analysis(sql string, v ...Value) (string, error) {

	return "", nil
}

func (s *SQL) Sql() string {
	build := strings.Builder{}
	var statement string
	if s.selects {
		build.WriteString("SELECT ")
		statement = strings.Join(s.column, ",")
		build.WriteString(statement)
	}

	if s.tables {
		build.WriteString(" FROM ")
		statement = strings.Join(s.table, ",")
		build.WriteString(statement)
	}

	if s.left {
		build.WriteString(" LEFT JOIN ")
		if s.o_n {
			build.WriteString("(")
		}
		statement = strings.Join(s.left_join, " ")
		build.WriteString(statement)
		if s.o_n {
			build.WriteString(")")
		}
	}

	if s.right {
		build.WriteString(" RIGHT JOIN ")
		if s.o_n {
			build.WriteString("(")
		}
		statement = strings.Join(s.right_join, " ")
		build.WriteString(statement)
		if s.o_n {
			build.WriteString(")")
		}
	}
	if s.o_n {
		build.WriteString(" ON (")
		statement = strings.Join(s.on, " ")
		build.WriteString(statement)
		build.WriteString(")")
	}

	if s.wheres {
		build.WriteString(" WHERE ")
		if s.condition {
			build.WriteString("(")
			s.where = append(s.where, ")")
		}
		statement = strings.Join(s.where, "")
		build.WriteString(statement)
	}

	return build.String()
}
