package utils

import (
	"fmt"
	"gitee.com/aurora-engine/aurora/utils/maputils"
	"gitee.com/aurora-engine/aurora/utils/sliceutils"
	"gitee.com/aurora-engine/aurora/utils/strutils"
	"math"
	"testing"
)

func TestSlice(t *testing.T) {
	arr := []int{1}
	t.Log(sliceutils.IsEmpty(arr))
	t.Log(sliceutils.StrSlice(1, 2.2, false))
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
	f := 13.36
	floor := fmt.Sprintf("%0.2f", math.Floor((f+0.05)*100)/100)
	t.Log(floor)
}
