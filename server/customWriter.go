package server

import "net/http"

type CustomWriter struct {
	http.ResponseWriter
	code int
	body []byte
}

func (w *CustomWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *CustomWriter) Write(body []byte) (int, error) {
	w.body = body
	return w.ResponseWriter.Write(body)
}
