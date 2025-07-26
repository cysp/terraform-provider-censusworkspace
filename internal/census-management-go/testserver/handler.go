package censusmanagementtestserver

import (
	cms "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go/server"
)

type handler struct {
	ts *CensusManagementTestServer
}

var _ cms.Handler = (*handler)(nil)
