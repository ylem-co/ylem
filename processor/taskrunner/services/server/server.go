package server

import (
	"net/http"
	"ylem_taskrunner/api"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Listen string
}

func (s *Server) Run() error {
	log.Info("Starting server listening on " + s.Listen)

	rtr := mux.NewRouter()

	rtr.HandleFunc("/hello/{name}/", api.ExampleHandler).Methods(http.MethodGet)
	http.Handle("/", rtr)

	return http.ListenAndServe(s.Listen, rtr)
}

func NewServer(listen string) *Server {
	s := &Server{
		Listen: listen,
	}

	return s
}
