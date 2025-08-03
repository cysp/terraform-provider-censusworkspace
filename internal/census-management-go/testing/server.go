package testing

import (
	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

type Server struct {
	server *cm.Server

	handler *Handler
}

func NewCensusManagementServer() (*Server, error) {
	handler := NewCensusManagementHandler()

	server, err := cm.NewServer(handler)
	if err != nil {
		return nil, err
	}

	return &Server{
		server:  server,
		handler: handler,
	}, nil
}

func (s *Server) Handler() *Handler {
	return s.handler
}
