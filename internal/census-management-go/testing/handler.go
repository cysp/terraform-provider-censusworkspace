package testing

import (
	"sync"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

type Handler struct {
	mu sync.Mutex

	Sources      map[string]*cm.SourceData
	sourceIDLast int64

	Models      multilevelMap2[string, string, *cm.SourceModelData]
	modelIDLast int64
}

var _ cm.Handler = (*Handler)(nil)

func NewCensusManagementHandler() *Handler {
	return &Handler{
		mu: sync.Mutex{},

		Sources: make(map[string]*cm.SourceData),
	}
}
