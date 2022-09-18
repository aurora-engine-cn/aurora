package timeutils

import (
	"gitee.com/aurora-engine/aurora/utils"
	"time"
)

const (
	datetime = "2006-01-02 03:04:05"
	date     = "2006-01-02"
)

// DateTime 获取当前时间
func DateTime(format ...string) string {
	f := datetime
	if format != nil {
		f = format[0]
	}
	return time.Now().Format(f)
}

// Date 获取当前日期
func Date(format ...string) string {
	f := date
	if format != nil {
		f = format[0]
	}
	return time.Now().Format(f)
}

// BeforeTime 获取 day 天之前的时间
func BeforeTime(day int) string {
	day = day * 24 * 60 * 60
	add := time.Now().Add(-time.Duration(day) * time.Second)
	return add.Format(datetime)
}

// AfterTime 获取 day 天之后的时间
func AfterTime(day int) string {
	day = day * 24 * 60 * 60
	add := time.Now().Add(-time.Duration(day) * time.Second)
	return add.Format(datetime)
}

// BeforeDate 获取 day 天之前的日期
func BeforeDate(day int) string {
	day = day * 24 * 60 * 60
	add := time.Now().Add(time.Duration(day) * time.Second)
	return add.Format(date)
}

// AfterDate 获取 day 天之后的日期
func AfterDate(day int) string {
	day = day * 24 * 60 * 60
	add := time.Now().Add(time.Duration(day) * time.Second)
	return add.Format(date)
}

func Time(v string) time.Time {
	parse, err := time.Parse(datetime, v)
	utils.UtilError(err)
	return parse
}
