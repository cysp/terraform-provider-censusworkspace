package testserver

import (
	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

type handler struct {
	ts *CensusManagementTestServer
}

var _ cm.Handler = (*handler)(nil)
