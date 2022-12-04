package utils

import (
	"fmt"
	"gitee.com/aurora-engine/aurora/utils/maputils"
	"gitee.com/aurora-engine/aurora/utils/sliceutils"
	"gitee.com/aurora-engine/aurora/utils/strutils"
	"gitee.com/aurora-engine/aurora/utils/uerr"
	"testing"
)

func TestSlice(t *testing.T) {
	arr := []int{1}
	//t.Log(sliceutils.IsEmpty(arr))
	//t.Log(sliceutils.StrSlice(1, 2.2, false))
	t.Log(sliceutils.Slice(arr))
}

func TestMap(t *testing.T) {
	t.Log(maputils.IsEmpty(map[string]float64{"1": 1}))
}

func TestString(t *testing.T) {
	t.Log(strutils.IsEmpty(""))
	t.Log(strutils.Int("1111"))

	t.Log(strutils.String(12))
	t.Log(strutils.String(12.235))
	t.Log(strutils.String(true))
}

func TestMath(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case error:
				fmt.Println(err.(error).Error())
			case uerr.UtilsError:
				fmt.Println(err.(uerr.UtilsError))
			}
		}
	}()
	panic(uerr.UtilsError("error"))
}
