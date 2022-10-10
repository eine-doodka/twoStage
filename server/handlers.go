package server

import "net/http"

type Handlers struct{}

func (h *Handlers) handleInit() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
func (h *Handlers) handleCommit() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
