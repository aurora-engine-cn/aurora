package route

import (
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
