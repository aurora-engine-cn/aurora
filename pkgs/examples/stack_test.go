package examples

import (
	"gitee.com/aurora-engine/aurora/pkgs/stack"
	"testing"
)

func TestStack(t *testing.T) {
	s := stack.Stack[any]{}
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Push(4)
	t.Log(s.Size())
	t.Log(s.IsEmpty())
	for !s.IsEmpty() {
		t.Log(s.Popup())
	}
	t.Log(s.Size())
	t.Log(s.IsEmpty())
}
