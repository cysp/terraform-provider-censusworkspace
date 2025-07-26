package testserver

import (
	"net/http/httptest"
	"sync"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

const (
	NonexistentID = "nonexistent"
)

type CensusManagementTestServer struct {
	mu *sync.Mutex

	httpTestServer *httptest.Server

	Sources      map[int64]*cm.SourceData
	sourceIDLast int64
}

func NewCensusManagementTestServer() *CensusManagementTestServer {
	ts := &CensusManagementTestServer{
		mu:      &sync.Mutex{},
		Sources: make(map[int64]*cm.SourceData),
	}

	h := &handler{ts: ts}

	server, err := cm.NewServer(h, h)
	if err != nil {
		return nil
	}

	ts.httpTestServer = httptest.NewServer(server)

	return ts
}

func (ts *CensusManagementTestServer) Server() *httptest.Server {
	return ts.httpTestServer
}
