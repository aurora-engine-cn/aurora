package uerr

import (
	"fmt"
	"strings"
)

// UtilError utils中的错误panic处理
func UtilError(err error, msg ...string) {
	if err != nil {
		errMsg := fmt.Sprintf("utils error: %s %s", strings.Join(msg, ""), err.Error())
		panic(errMsg)
	}
}
