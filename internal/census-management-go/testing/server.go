package testing

import (
	"net/http"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

type Server struct {
	server *cm.Server

	h   *Handler
	sec *SecurityHandler
}

var _ http.Handler = (*Server)(nil)

func NewCensusManagementServer() (*Server, error) {
	h := NewCensusManagementHandler()

	sec := NewCensusManagementSecurityHandler()

	server, err := cm.NewServer(h, sec)
	if err != nil {
		return nil, err
	}

	return &Server{
		server: server,
		h:      h,
		sec:    sec,
	}, nil
}

func (s *Server) Handler() *Handler {
	return s.h
}

func (s *Server) SecurityHandler() *SecurityHandler {
	return s.sec
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.server.ServeHTTP(w, r)
}
