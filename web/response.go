package web

import (
	"net/http"
)

type Response interface {
	SetHeader(string, string)
	AddHeader(string, string)
}

type Rew struct {
	http.ResponseWriter
}
