package uerr

import (
	"fmt"
	"strings"
)

// UtilsError 用于 aurora 包下的 utils panic 错误 以便第三方自定义 错误捕捉器能够分类处理
type UtilsError string

// UtilError utils中的错误panic处理
// msg 附加消息，会放在 具体error消息之前
func UtilError(err error, msg ...string) {
	if err != nil {
		errMsg := fmt.Sprintf("utils msg: %s. err: %s", strings.Join(msg, ""), err.Error())
		panic(UtilsError(errMsg))
	}
}
