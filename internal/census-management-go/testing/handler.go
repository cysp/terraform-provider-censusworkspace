package testing

import (
	"sync"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

type Handler struct {
	mu sync.Mutex
}

var _ cm.Handler = (*Handler)(nil)

func NewCensusManagementHandler() *Handler {
	return &Handler{
		mu: sync.Mutex{},
	}
}
