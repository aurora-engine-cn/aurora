package sqlbuild

import "strings"

/*
	处理 转移
*/

const ()

func replaceChart(s ...string) []string {
	column := make([]string, 0)
	c := ""
	for i := 0; i < len(s); i++ {
		str := s[i]
		if str == "*" {
			return s
		}
		var arr []string
		if strings.Contains(str, ",") {
			arr = strings.Split(str, ",")
			for j := 0; j < len(arr); j++ {
				c = replace(arr[j])
				column = append(column, c)
			}
			continue
		}
		c = replace(str)
		column = append(column, c)
	}
	return column
}

// 处理 select 选择字段
func selectReplaceChart(s ...string) []string {
	column := make([]string, 0)
	c := ""
	for i := 0; i < len(s); i++ {
		str := s[i]
		if str == "*" {
			return s
		}
		var arr []string
		if strings.Contains(str, ",") {
			arr = strings.Split(str, ",")
			for j := 0; j < len(arr); j++ {
				c = replaceSelect(arr[j])
				column = append(column, c)
			}
			continue
		}
		c = replaceSelect(str)
		column = append(column, c)
	}
	return column
}

// 处理 where 条件字段
func whereReplaceChart(s ...string) []string {
	column := make([]string, 0)
	c := ""
	for i := 0; i < len(s); i++ {
		str := s[i]
		var arr []string
		w := make([]string, 0)
		// 检查 是否拼接条件
		if strings.Contains(str, AND) {
			arr = strings.Split(str, AND)
			for j := 0; j < len(arr); j++ {
				// 校验 子句 是否存在条件 存在组合条件就继续递归 通过 AND 分割的 情况下 只需要校验 or 和 OR 情况
				if strings.Contains(arr[j], or) {
					arr2 := strings.Split(arr[j], or)
					w2 := whereReplaceChart(arr2...)
					w = append(w, strings.Join(w2, or))
					continue
				}
				if strings.Contains(arr[j], OR) {
					arr2 := strings.Split(arr[j], OR)
					w2 := whereReplaceChart(arr2...)
					w = append(w, strings.Join(w2, OR))
					continue
				}
				c = replaceWhere(arr[j])
				w = append(w, c)
			}
			// strings.Join(w, AND) 如何分解的条件 就如何拼接起来 对应上面的 arr = strings.Split(str, AND)
			column = append(column, strings.Join(w, AND))
			continue
		}
		if strings.Contains(str, and) {
			arr = strings.Split(str, and)
			for j := 0; j < len(arr); j++ {
				if strings.Contains(arr[j], or) {
					arr2 := strings.Split(arr[j], or)
					w2 := whereReplaceChart(arr2...)
					w = append(w, strings.Join(w2, or))
					continue
				}
				if strings.Contains(arr[j], OR) {
					arr2 := strings.Split(arr[j], OR)
					w2 := whereReplaceChart(arr2...)
					w = append(w, strings.Join(w2, OR))
					continue
				}
				c = replaceWhere(arr[j])
				w = append(w, c)
			}
			column = append(column, strings.Join(w, and))
			continue
		}

		if strings.Contains(str, OR) {
			arr = strings.Split(str, OR)
			for j := 0; j < len(arr); j++ {
				if strings.Contains(arr[j], and) {
					arr2 := strings.Split(arr[j], and)
					w2 := whereReplaceChart(arr2...)
					w = append(w, strings.Join(w2, and))
					continue
				}
				if strings.Contains(arr[j], AND) {
					arr2 := strings.Split(arr[i], AND)
					w2 := whereReplaceChart(arr2...)
					w = append(w, strings.Join(w2, AND))
					continue
				}
				c = replaceWhere(arr[j])
				w = append(w, c)
			}
			column = append(column, strings.Join(w, OR))
			continue
		}
		if strings.Contains(str, or) {
			arr = strings.Split(str, or)
			for j := 0; j < len(arr); j++ {
				if strings.Contains(arr[j], and) {
					arr2 := strings.Split(arr[j], and)
					w2 := whereReplaceChart(arr2...)
					w = append(w, strings.Join(w2, and))
					continue
				}
				if strings.Contains(arr[j], AND) {
					arr2 := strings.Split(arr[j], AND)
					w2 := whereReplaceChart(arr2...)
					w = append(w, strings.Join(w2, AND))
					continue
				}
				c = replaceWhere(arr[j])
				w = append(w, c)
			}
			column = append(column, strings.Join(w, or))
			continue
		}
		c = replaceWhere(str)
		column = append(column, c)
	}
	return column
}

func replaceWhere(s string) string {
	if strings.Contains(s, "=") {
		arr := strings.Split(s, "=")
		return replace(arr[0]) + "=" + arr[1]
	}
	return ""
}

func replaceSet(s string) string {
	if strings.Contains(s, "=") {
		arr := strings.Split(s, "=")
		return replace(arr[0]) + "=" + arr[1]
	}
	return ""
}

func replaceSelect(s string) string {
	//检查选择字段别名处理
	if strings.Contains(s, AS) {
		arr := strings.Split(s, AS)
		s2 := replace(arr[0])
		s3 := replace(arr[1])
		return strings.Join([]string{s2, s3}, AS)
	}
	if strings.Contains(s, as) {
		arr := strings.Split(s, as)
		s2 := replace(arr[0])
		s3 := replace(arr[1])
		return strings.Join([]string{s2, s3}, as)
	}
	// 待添加对 空格别名的支持
	//if strings.Contains(s, " ") {
	//	arr := strings.Split(s, " ")
	//	s2 := replace(arr[0])
	//	s3 := replace(arr[1])
	//	return strings.Join([]string{s2, s3}, " ")
	//}
	return replace(s)
}

// 处理 sql 选择字段 条件字段 和更新字段的 关键字区别
func replace(s string) string {
	if strings.Contains(s, ".") {
		arr := strings.Split(s, ".")
		v := ""
		for i := 0; i < len(arr); i++ {
			// 去掉多余空格以防添加符号导致 sql 语句出错
			arr[i] = strings.TrimSpace(arr[i])
			v += chart + arr[i] + chart + "."
		}
		return v[:len(v)-1]
	}
	s = strings.TrimSpace(s)
	return chart + s + chart
}
