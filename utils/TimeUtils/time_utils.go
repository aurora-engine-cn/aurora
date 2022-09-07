package TimeUtils

import (
	"gitee.com/aurora-engine/aurora/utils"
	"time"
)

const (
	formatYYYYMMDDHHMMSS = "2006-01-02 03:04:05"
	formatYYYYMMDD       = "2006-01-02"
)

// DateTime 获取当前时间
func DateTime() string {
	return time.Now().Format(formatYYYYMMDDHHMMSS)
}

// Date 获取当前日期
func Date() string {
	return time.Now().Format(formatYYYYMMDD)
}

func Time(v string) time.Time {
	parse, err := time.Parse(formatYYYYMMDDHHMMSS, v)
	utils.UtilError(err)
	return parse
}
