package example

import (
	"fmt"
	"gitee.com/aurora-engine/aurora"
	"testing"
)

type Aaa struct {
	Name string
	*Bbb
}

type Bbb struct {
	Name string
	*Aaa
	*Ccc
}

type Ccc struct {
	Name string
}

func TestIoc(t *testing.T) {
}

type TestServer struct {
	*Aaa
	*Bbb
	Ccc `ref:"a"`
}

func Before(ctx aurora.Ctx) bool {
	fmt.Println("Before")
	return true
}

func (s *TestServer) GetName() string {
	return s.Aaa.Bbb.Name
}

func (s *TestServer) GetUpdate() {
	s.Aaa.Name = "b"
}
