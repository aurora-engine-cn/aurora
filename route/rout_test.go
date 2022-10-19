package route

import (
	"strings"
	"testing"
)

func BenchmarkRESTFul(b *testing.B) {
	n := &node{
		Path:       "",
		FullPath:   "/{name}/{2}",
		Count:      0,
		middleware: nil,
		Control:    nil,
		Child:      nil,
	}
	for i := 0; i < b.N; i++ {
		RESTFul(n, "/a/222")
	}
}
func BenchmarkRESTFul2(b *testing.B) {
	n := &node{
		Path:       "",
		FullPath:   "/{name}/{2}",
		Count:      0,
		middleware: nil,
		Control:    nil,
		Child:      nil,
	}
	for i := 0; i < b.N; i++ {
		analysisRESTFul(n, "/a/222")
	}
}

func BenchmarkCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.Count("/a/222/sad/asc/a//das/dasf/asdfaf/as/f/asf/as/f", "/")
	}
}

func BenchmarkCountChar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CountChar("/a/222/sad/asc/a//das/dasf/asdfaf/as/f/asf/as/f", "/")
	}
}

func CountChar(parent, sub string) int {
	c := 0
	for i := 0; i < len(parent); i++ {
		if parent[i:i+1] == sub {
			c++
		}
	}
	return c
}
