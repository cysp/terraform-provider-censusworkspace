package censusmanagementtestserver

import (
	"net/http"
	"net/http/httptest"
	"sync"

	cm "github.com/cysp/terraform-provider-census/internal/census-management-go"
)

const (
	NonexistentID = "nonexistent"
)

type CensusManagementTestServer struct {
	mu *sync.Mutex

	httpTestServer *httptest.Server
	serveMux       *http.ServeMux

	sources map[int64]*cm.Source
}

func NewCensusManagementTestServer() *CensusManagementTestServer {
	testserver := &CensusManagementTestServer{
		mu:      &sync.Mutex{},
		sources: make(map[int64]*cm.Source),
	}

	testserver.serveMux = http.NewServeMux()
	testserver.httpTestServer = httptest.NewServer(testserver.serveMux)

	testserver.setupUserHandler()
	testserver.setupPersonalAccessTokenHandlers()
	testserver.setupOrganizationAppDefinitionHandlers()
	testserver.setupOrganizationAppDefinitionResourceProviderHandlers()
	testserver.setupOrganizationAppDefinitionResourceTypeHandlers()
	testserver.setupSpaceEnablementsHandlers()
	testserver.setupSpaceAPIKeyHandlers()
	testserver.SetupSpaceEnvironmentAppInstallationHandlers()
	testserver.setupSpaceEnvironmentContentTypeHandlers()
	testserver.SetupSpaceEnvironmentExtensionHandlers()
	testserver.setupSpacePreviewAPIKeyHandlers()
	testserver.setupSpaceRoleHandlers()
	testserver.setupSpaceWebhookDefinitionHandlers()

	return testserver
}

func (ts *CensusManagementTestServer) Server() *httptest.Server {
	return ts.httpTestServer
}
