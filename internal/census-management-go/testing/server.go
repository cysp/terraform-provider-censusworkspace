package testing

import (
	"sync"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

type Server struct {
	server *cm.Server

	h   *Handler
	sec *SecurityHandler
}

func NewCensusManagementServer() (*Server, error) {
	h := &Handler{
		mu: sync.Mutex{},

		Sources: make(map[string]*cm.SourceData),
	}

	sec := &SecurityHandler{
		mu: sync.Mutex{},
	}

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

func (s *Server) Server() *cm.Server {
	return s.server
}

func (s *Server) Handler() *Handler {
	return s.h
}

func (s *Server) SecurityHandler() *SecurityHandler {
	return s.sec
}
