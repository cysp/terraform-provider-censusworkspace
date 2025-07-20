package testing

import (
	"net/http"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

type Server struct {
	server *cm.Server

	handler         *Handler
	securityHandler *SecurityHandler
}

var _ http.Handler = (*Server)(nil)

func NewCensusManagementServer() (*Server, error) {
	handler := NewCensusManagementHandler()

	securityHandler := NewCensusManagementSecurityHandler()

	server, err := cm.NewServer(handler, securityHandler)
	if err != nil {
		return nil, err
	}

	return &Server{
		server:          server,
		handler:         handler,
		securityHandler: securityHandler,
	}, nil
}

func (s *Server) Handler() *Handler {
	return s.handler
}

func (s *Server) SecurityHandler() *SecurityHandler {
	return s.securityHandler
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.server.ServeHTTP(w, r)
}
