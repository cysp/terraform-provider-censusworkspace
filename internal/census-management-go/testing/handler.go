package testing

import (
	"sync"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

type Handler struct {
	mu sync.Mutex

	Destinations      map[string]*cm.DestinationData
	destinationIDLast int64

	Sources      map[string]*cm.SourceData
	sourceIDLast int64
}

var _ cm.Handler = (*Handler)(nil)

func NewCensusManagementHandler() *Handler {
	return &Handler{
		mu: sync.Mutex{},

		Destinations: make(map[string]*cm.DestinationData),

		Sources: make(map[string]*cm.SourceData),
	}
}
