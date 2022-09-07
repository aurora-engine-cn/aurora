package examples

import (
	"fmt"
	"gitee.com/aurora-engine/pkgs/list"
	"testing"
)

func TestLinkList_Add(t *testing.T) {
	l := list.LinkList[int]{}
	l.Add(0)
	l.Add(1)
	fmt.Println(l.Get(0))
	fmt.Println(l.Get(1))
	fmt.Println(l.Get(2))
}
