package examples

import (
	"fmt"
	"gitee.com/aurora-engine/aurora/pkgs/list"
	"testing"
)

func TestList(t *testing.T) {
	l := list.ArrayList[int]{1, 2, 3}
	for _, v := range l {
		fmt.Println(v)
	}
	l.Remove(0)
	t.Log(l)
}
