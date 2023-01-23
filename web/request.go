package web

import "net/http"

type Request interface {
	GetHeader(string) string

	GetString(string) string
	GetInt(string) int
	GetInt64(string) int64
	GetFloat64(string) float64
	GetFloat32(string) float32
	GetBool(string) bool

	GetStrings(string) []string
	GetInts(string) []int
	GetInt64s(string) []int64
	GetFloat64s(string) []float64
	GetFloat32s(string) []float32
	GetBools(string) []bool

	Post() any
}

type Req struct {
	http.Request
}
