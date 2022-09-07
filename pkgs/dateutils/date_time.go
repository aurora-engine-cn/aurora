package dateutils

import "time"

/*
	日期时间 工具
*/

const (
	datetime = "2006-01-02 03:04:05"
	date     = "2006-01-02"
)

// DateTime 获取时间字符串
func DateTime(format ...string) string {
	f := datetime
	if format != nil {
		f = format[0]
	}
	return time.Now().Format(f)
}

// Date 获取日期字符串
func Date(format ...string) string {
	f := date
	if format != nil {
		f = format[0]
	}
	return time.Now().Format(f)
}

// BeforeTime 获取 day 之前的时间
func BeforeTime(day int) string {
	day = day * 24 * 60 * 60
	add := time.Now().Add(-time.Duration(day) * time.Second)
	return add.Format(datetime)
}

// BeforeDate 获取 day 之前的日期
func BeforeDate(day int) string {
	day = day * 24 * 60 * 60
	add := time.Now().Add(-time.Duration(day) * time.Second)
	return add.Format(date)
}
