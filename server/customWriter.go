package server

import "net/http"

type CustomWriter struct {
	http.ResponseWriter
	code int
}

func (w CustomWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}
