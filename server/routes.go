package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	ctxkeyReqId = iota
)

type Server struct {
	router *mux.Router
	log    *logrus.Logger
	handle *Handlers
}

func NewServer(handlers *Handlers) *Server {
	s := &Server{
		router: mux.NewRouter(),
		log:    logrus.New(),
		handle: handlers,
	}
	s.configureRouter()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) configureRouter() {
	s.router.Use(s.loggingMw)
	s.router.Handle("/init", s.handle.handleInit()).Methods("GET")
	s.router.Handle("/commit", s.handle.handleCommit()).Methods("POST")
	s.router.PathPrefix("/").Handler(s.defaultAnswer())
}

func (s *Server) loggingMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := uuid.New()
		ctx := context.WithValue(context.Background(),
			ctxkeyReqId,
			reqId,
		)
		logger := s.log.WithFields(
			map[string]interface{}{
				"remote_addr": r.RemoteAddr,
				"request_id":  reqId,
			})
		wr := &CustomWriter{w, http.StatusOK, []byte{}}
		next.ServeHTTP(wr, r.WithContext(ctx))
		logger.Infof("Response with code %v: %v", wr.code, http.StatusText(wr.code))
		if wr.code >= 500 {
			logger.Warnf("Server-side error with body: %v", string(wr.body))
		}
	})
}

func (s *Server) defaultAnswer() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	})
}
